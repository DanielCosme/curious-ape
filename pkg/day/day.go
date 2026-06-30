package day

import (
	"errors"
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/persistence"
	"github.com/stephenafamo/bob"
)

type App struct {
	daysDB   core.DayRepository
	habitsDB core.HabitRepository
}

func New(db bob.DB) *App {
	return &App{
		daysDB:   persistence.NewDays(db),
		habitsDB: persistence.NewHabits(db),
	}
}

func (a *App) Month(date core.Date, order core.OrderParam) ([]core.Day, error) {
	if time.Now().Month() != date.Time().Month() {
		date = date.LastDayOfTheMonth()
	}

	day, err := a.daysDB.Get(core.DayParams{Date: date})
	if core.IfErrNNotFound(err) {
		return nil, err
	}
	daysOfTheMonth := date.RangeMonth()
	if day.IsZero() {
		for _, date := range daysOfTheMonth {
			if _, err := a.GetOrCreate(date); err != nil {
				return nil, err
			}
		}
	}

	return a.daysDB.Find(core.DayParams{Dates: daysOfTheMonth, Order: order})
}

func (a *App) GetOrCreate(d core.Date) (day core.Day, err error) {
	day, err = a.daysDB.Get(core.DayParams{Date: d})
	if core.IfErrNNotFound(err) {
		return
	}
	if day.IsZero() {
		day, err = a.daysDB.Create(d)
		if err != nil {
			return
		}

		// I want this upserts to be Events that the Habits-Thing Listens to...
		// But I want this
		_, e1 := a.habitsDB.Upsert(core.Habit{Date: day.Date, State: core.HabitStateNoInfo, Type: core.HabitTypeWakeUp})
		_, e2 := a.habitsDB.Upsert(core.Habit{Date: day.Date, State: core.HabitStateNoInfo, Type: core.HabitTypeFitness})
		_, e3 := a.habitsDB.Upsert(core.Habit{Date: day.Date, State: core.HabitStateNoInfo, Type: core.HabitTypeDeepWork})
		_, e4 := a.habitsDB.Upsert(core.Habit{Date: day.Date, State: core.HabitStateNoInfo, Type: core.HabitTypeEatHealthy})
		err = errors.Join(e1, e2, e3, e4)
		if err != nil {
			return
		}
		return a.daysDB.Get(core.DayParams{ID: day.ID})
	}
	return
}
