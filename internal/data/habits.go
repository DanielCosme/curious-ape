package data

import (
	"database/sql"
	"errors"

	"github.com/danielcosme/curious-ape/internal/validator"
)

type Habit struct {
	ID     int    `json:"id"`
	State  string `json:"state"`
	Date   string `json:"date"`
	Origin string `json:"origin"`
	Type   string `json:"type"`
}

type HabitModel struct {
	DB *sql.DB
}

func (hm *HabitModel) GetAll() ([]*Habit, error) {
	query := `
		SELECT id, state, date, origin, type FROM habits`
	rows, err := hm.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	habits := []*Habit{}
	for rows.Next() {
		var habit Habit

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

		habits = append(habits, &habit)
	}

	return habits, nil
}

func (fh *HabitModel) Insert(habit *Habit) error {
	query := `
		INSERT INTO habits (state, "date", "type", origin)
		VALUES ($1, $2, $3, $4)
		RETURNING id `
	args := []interface{}{habit.State, habit.Date, habit.Type, habit.Origin}
	return fh.DB.QueryRow(query, args...).Scan(&habit.ID)
}

func (fh *HabitModel) Get(id int) (*Habit, error) {
	query := `
		SELECT id, state, "date", "type", origin FROM habits
		WHERE id = $1`
	var habit Habit
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
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &habit, nil
}

func (fh *HabitModel) UpdateOrCreate(habit *Habit) error {
	stm := `SELECT id FROM habits WHERE "date" = $1 and "type" = $2`
	q, _ := fh.DB.Query(stm, habit.Date, habit.Type)
	if !q.Next() {
		if q.Err() == nil {
			return fh.Insert(habit)
		}
		return q.Err()
	}

	if err := q.Scan(&habit.ID); err != nil {
		return err
	}

	if err := fh.Update(habit); err != nil {
		return err
	}

	return nil
}

func (fh *HabitModel) Update(habit *Habit) error {
	stm := `
		UPDATE habits
		SET state = $1, "date" = $2, "type" = $3, "origin" = $4
		WHERE id = $5
		RETURNING state, "date", "type"`

	args := []interface{}{
		habit.State,
		habit.Date,
		habit.Type,
		habit.Origin,
		habit.ID,
	}

	return fh.DB.QueryRow(stm, args...).Scan(&habit.State, &habit.Date, &habit.Type)
}

func (fh *HabitModel) UpdateByDate(habit *Habit) error {
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
		return ErrRecordNotFound
	}

	return nil
}

// TODO Add more robust validation
func ValidateHabit(v *validator.Validator, habit *Habit) {
	v.Check(habit.Date != "", "date", "must be provided")
	v.Check(len([]rune(habit.Date)) == 10, "date", "must be exactly 10 characters long")
}
