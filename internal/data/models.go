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
	Habits       HabitModel
	Users        UserModel
	SleepRecords SleepRecordModel
	Tokens       AuthTokenModel
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Habits:       HabitModel{DB: db},
		Users:        UserModel{DB: db},
		SleepRecords: SleepRecordModel{DB: db},
		Tokens:       AuthTokenModel{DB: db},
	}
}
