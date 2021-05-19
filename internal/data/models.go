package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	FoodHabits FoodHabitModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		FoodHabits: FoodHabitModel{DB: db},
	}
}
