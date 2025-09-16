package persistence

import (
	"github.com/stephenafamo/bob"
)

type Database struct {
	Users    Users
	Days     Days
	Habits   Habits
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
