package models

import (
	"database/sql"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/models/pg"
)

type DB struct {
	Habits         core.HabitModel
	Users          core.UserModel
	SleepRecords   core.SleepRecordModel
	Tokens         core.AuthTokenModel
	WorkRecords    core.WorkModel
	FitnessRecords core.FitnessModel
}

func NewModels(db *sql.DB) *DB {
	return &DB{
		Habits:         &pg.HabitModel{DB: db},
		Users:          &pg.UserModel{DB: db},
		SleepRecords:   &pg.SleepRecordModel{DB: db},
		Tokens:         &pg.AuthTokenModel{DB: db},
		WorkRecords:    &pg.WorkModel{DB: db},
		FitnessRecords: &pg.FitnessModel{DB: db},
	}
}
