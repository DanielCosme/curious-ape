package persistence

import (
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/stephenafamo/bob"
)

type Database struct {
	Users    Users
	Days     core.DayRepository
	Habits   core.HabitRepository
	Sleep    SleepLogs
	Fitness  FitnessLogs
	DeepWork DeepWorkLogs
	Auths    Auths
	executor bob.DB
}

func New(executor bob.DB) *Database {
	return &Database{
		Users:    Users{db: executor},
		Days:     Days{db: executor},
		Habits:   Habits{db: executor},
		Sleep:    SleepLogs{db: executor},
		Fitness:  FitnessLogs{db: executor},
		DeepWork: DeepWorkLogs{db: executor},
		Auths:    Auths{db: executor},
		executor: executor,
	}
}
