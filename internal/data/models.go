package data

import (
	"database/sql"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/data/pg"
)

type Models struct {
	Habits         core.HabitModel
	Users          UserModel
	SleepRecords   SleepRecordModel
	Tokens         AuthTokenModel
	WorkRecords    WorkModel
	FitnessRecords core.FitnessModel
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Habits:         &pg.HabitModel{DB: db},
		Users:          UserModel{DB: db},
		SleepRecords:   SleepRecordModel{DB: db},
		Tokens:         AuthTokenModel{DB: db},
		WorkRecords:    WorkModel{DB: db},
		FitnessRecords: &pg.FitnessModel{DB: db},
	}
}
