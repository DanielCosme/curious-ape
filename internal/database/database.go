package database

import (
	"github.com/aarondl/opt/omit"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database/gen/models"
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

func (d *Database) HabitGetOrCreate(date core.Date, habitType core.HabitType) (core.Habit, error) {
	var res core.Habit
	hc, err := HabitCategoryParams{Type: habitType}.BuildQuery(d.executor).One()
	if err != nil {
		return res, catchErr("habit get or create", err)
	}

	day, err := d.Days.GetOrCreate(DayParams{Date: date})
	if err != nil {
		return res, catchErr("habit get or create (query day)", err)
	}

	res, err = d.Habits.Get(HabitParams{DayID: day.ID, CategoryID: hc.ID})
	if IfNotFoundErr(err) {
		return res, err
	}
	if res.IsZero() {
		// Create new habit.
		return d.Habits.Create(models.HabitSetter{
			DayID:           omit.From(day.ID),
			HabitCategoryID: omit.From(hc.ID),
			State:           omit.From(string(core.HabitStateNoInfo)),
		})
	}
	return res, nil
}
