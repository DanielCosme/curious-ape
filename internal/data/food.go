package data

import (
	"database/sql"
	"errors"

	"github.com/danielcosme/curious-ape/internal/validator"
	"github.com/lib/pq"
)

type FoodHabit struct {
	ID    int      `json:"-"` // The - directive hides it from json
	State bool     `json:"state"`
	Date  string   `json:"date"`
	Tags  []string `json:"tags,omitempty"` // Tags: 16/8_fast, lion, calorie_deficit, calorie_surplus.
}

type FoodHabitModel struct {
	DB *sql.DB
}

func (fh *FoodHabitModel) Insert(habit *FoodHabit) error {
	query := `
		INSERT INTO food_habits (state, "date", tags)
		VALUES ($1, $2, $3)
		RETURNING id `
	args := []interface{}{habit.State, habit.Date, pq.Array(habit.Tags)}
	return fh.DB.QueryRow(query, args...).Scan(&habit.ID)
}

func (fh *FoodHabitModel) Get(date string) (*FoodHabit, error) {
	query := `
		SELECT id, state, date, tags FROM food_habits 
		WHERE "date" = $1`
	var habit FoodHabit
	err := fh.DB.QueryRow(query, date).Scan(
		&habit.ID,
		&habit.State,
		&habit.Date,
		pq.Array(&habit.Tags),
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

func (fh *FoodHabitModel) Update(habit *FoodHabit) error {
	stm := `
		UPDATE food_habits
		SET state = $1, "date" = $2, tags = $3
		WHERE "date" = $2
		RETURNING state, "date", tags`

	args := []interface{}{
		habit.State,
		habit.Date,
		pq.Array(habit.Tags),
	}

	return fh.DB.QueryRow(stm, args...).Scan(&habit.State, &habit.Date, pq.Array(&habit.Tags))
}

func (fh *FoodHabitModel) Delete(id int) error {
	query := `
		DELETE FROM food_habits
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

func ValidateFoodHabit(v *validator.Validator, habit *FoodHabit) {
	v = validator.New()
	v.Check(habit.Date != "", "date", "must be provided")
	v.Check(len([]rune(habit.Date)) == 10, "date", "must be exactly 10 characters long")

	v.Check(len(habit.Tags) < 5, "tags", "must not have more than 5 tags")
	v.Check(validator.Unique(habit.Tags), "tags", "must not have duplicate values")
}
