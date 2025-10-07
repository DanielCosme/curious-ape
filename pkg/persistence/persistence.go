package persistence

import (
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/stephenafamo/bob"
)

type Database struct {
	Users    Users
	Days     core.DayRepository
	Habits   core.HabitRepository
	Sleep    core.SleepLogRepository
	Fitness  core.FitnessLogRepository
	DeepWork core.DeepWorkLogRepository
	Auths    Auths
	executor bob.DB
}

func New(executor bob.DB) *Database {
	return &Database{
		Users:    Users{db: executor},
		Days:     &Days{db: executor},
		Habits:   &Habits{db: executor},
		Sleep:    &SleepLogs{db: executor},
		Fitness:  &FitnessLogs{db: executor},
		DeepWork: &DeepWorkLogs{db: executor},
		Auths:    Auths{db: executor},
		executor: executor,
	}
}
