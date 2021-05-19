package data

import (
	"database/sql"

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
		Insert into food_habits (state, "date", tags)
		values ($1, $2, $3)
		Returning id `
	args := []interface{}{habit.State, habit.Date, pq.Array(habit.Tags)}
	return fh.DB.QueryRow(query, args...).Scan(&habit.ID)
}

func (fh *FoodHabitModel) Get(id int) (*FoodHabit, error) {
	return nil, nil
}

func (fh *FoodHabitModel) Update(habit *FoodHabit) error {
	return nil
}

func (fh *FoodHabitModel) Delete(id int) error {
	return nil
}

func ValidateFoodHabit(v *validator.Validator, habit *FoodHabit) {
	v = validator.New()
	v.Check(habit.Date != "", "date", "must be provided")
	v.Check(len([]rune(habit.Date)) == 10, "date", "must be exactly 10 characters long")

	v.Check(len(habit.Tags) < 5, "tags", "must not have more than 5 tags")
	v.Check(validator.Unique(habit.Tags), "tags", "must not have duplicate values")
}
