package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	FoodHabits   FoodHabitModel
	Users        UserModel
	SleepRecords SleepRecordModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		FoodHabits:   FoodHabitModel{DB: db},
		Users:        UserModel{DB: db},
		SleepRecords: SleepRecordModel{DB: db},
	}
}
