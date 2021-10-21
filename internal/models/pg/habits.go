package pg

import (
	"context"
	"database/sql"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/errors"
	"regexp"
	"time"

	"github.com/danielcosme/curious-ape/internal/validator"
)

type HabitModel struct {
	DB *sql.DB
}

func (hm *HabitModel) GetAll() ([]core.Habit, error) {
	query := `
		SELECT id, state, date, origin, type FROM habits`
	rows, err := hm.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	habits := []core.Habit{}
	for rows.Next() {
		var habit core.Habit

		err := rows.Scan(
			&habit.ID,
			&habit.State,
			&habit.Date,
			&habit.Origin,
			&habit.Type,
		)
		if err != nil {
			return nil, err
		}

		habits = append(habits, habit)
	}

	return habits, nil
}

func (fh *HabitModel) Insert(habit *core.Habit) error {
	query := `
		INSERT INTO habits (state, "date", "type", origin)
		VALUES ($1, $2, $3, $4)
		RETURNING id `
	args := []interface{}{habit.State, habit.Date, habit.Type, habit.Origin}
	return fh.DB.QueryRow(query, args...).Scan(&habit.ID)
}

func (fh *HabitModel) Get(id int) (*core.Habit, error) {
	query := `
		SELECT id, state, "date", "type", origin FROM habits
		WHERE id = $1`
	var habit core.Habit
	err := fh.DB.QueryRow(query, id).Scan(
		&habit.ID,
		&habit.State,
		&habit.Date,
		&habit.Type,
		&habit.Origin,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &habit, nil
}

func (fh *HabitModel) UpdateOrCreate(habit *core.Habit) error {
	stm := `SELECT id FROM habits WHERE "date" = $1 and "type" = $2`
	q, _ := fh.DB.Query(stm, habit.Date, habit.Type)
	if !q.Next() {
		if q.Err() == nil {
			return fh.Insert(habit)
		}
		return q.Err()
	}
	defer q.Close()

	if err := q.Scan(&habit.ID); err != nil {
		return err
	}

	err := fh.Update(habit)
	if err != nil {
		return err
	}

	return nil
}

func (fh *HabitModel) Update(h *core.Habit) error {
	stm := `
		UPDATE habits
		SET state = $2, "origin" = $3
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err := fh.DB.ExecContext(ctx, stm, &h.ID, &h.State, &h.Origin)
	if err != nil {
		return err
	}
	return nil
}

func (fh *HabitModel) UpdateByDate(habit *core.Habit) error {
	stm := `
		UPDATE habits
		SET state = $1, origin = $2
		WHERE "date" = $3 AND "type" = $4
		RETURNING state, "date", "type"`

	args := []interface{}{
		habit.State,
		habit.Origin,
		habit.Date,
		habit.Type,
	}

	return fh.DB.QueryRow(stm, args...).Scan(&habit.State, &habit.Date, &habit.Type)
}

func (fh *HabitModel) Delete(id int) error {
	query := `
		DELETE FROM habits
		WHERE id = $1`
	result, err := fh.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.ErrRecordNotFound
	}

	return nil
}

// TODO Add more robust validation
// validation for state and type + robust valiation for date
func ValidateHabit(v *validator.Validator, habit *core.Habit) {
	v.Check(habit.Date != "", "date", "must be provided")
	v.Check(len([]rune(habit.Date)) == 10, "date", "must be exactly 10 characters long")
	state := validator.Matches(habit.State, regexp.MustCompile(`^(yes|no|no_info)$`))
	v.Check(state, "state", "state must be yes, no or no_info")

	typeOfHabit := validator.Matches(habit.Type, regexp.MustCompile(`^(sleep|food|work|fitness)$`))
	v.Check(typeOfHabit, "type", "type must be sleep/food/work/fitness")
}
