package persistence

import (
	"danicos.dev/daniel/curious-ape/pkg/core"
	"github.com/stephenafamo/bob"
)

type Database struct {
	Users     Users
	Days      core.DayRepository
	Habits    core.HabitRepository
	Sleep     core.SleepLogRepository
	Fitness   core.FitnessLogRepository
	DeepWork  core.DeepWorkLogRepository
	Deadlines core.DeadlineRepository
	Auths     Auths
	executor  bob.DB
}

func New(executor bob.DB) *Database {
	return &Database{
		Users:     Users{db: executor},
		Days:      NewDays(executor),
		Habits:    NewHabits(executor),
		Sleep:     &SleepLogs{db: executor},
		Fitness:   &FitnessLogs{db: executor},
		DeepWork:  &DeepWorkLogs{db: executor},
		Auths:     Auths{db: executor},
		Deadlines: &Deadlines{db: executor},
		executor:  executor,
	}
}
