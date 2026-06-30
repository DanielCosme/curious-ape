package day

import (
	"errors"
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
)

func Month(date core.Date, order core.OrderParam) ([]core.Day, error) {
	if time.Now().Month() != date.Time().Month() {
		date = date.LastDayOfTheMonth()
	}

	day, err := daysDB.Get(core.DayParams{Date: date})
	if core.IfErrNNotFound(err) {
		return nil, err
	}
	daysOfTheMonth := date.RangeMonth()
	if day.IsZero() {
		for _, date := range daysOfTheMonth {
			if _, err := GetOrCreate(date); err != nil {
				return nil, err
			}
		}
	}

	return daysDB.Find(core.DayParams{Dates: daysOfTheMonth, Order: order})
}

func GetOrCreate(d core.Date) (day core.Day, err error) {
	day, err = daysDB.Get(core.DayParams{Date: d})
	if core.IfErrNNotFound(err) {
		return
	}
	if day.IsZero() {
		day, err = daysDB.Create(d)
		if err != nil {
			return
		}
		// TODOS: Implement the repository with SQLC

		// I want this upserts to be Events that the Habits-Thing Listens to...
		// But I want this
		_, e1 := habitsDB.Upsert(core.Habit{Date: day.Date, State: core.HabitStateNoInfo, Type: core.HabitTypeWakeUp})
		_, e2 := habitsDB.Upsert(core.Habit{Date: day.Date, State: core.HabitStateNoInfo, Type: core.HabitTypeFitness})
		_, e3 := habitsDB.Upsert(core.Habit{Date: day.Date, State: core.HabitStateNoInfo, Type: core.HabitTypeDeepWork})
		_, e4 := habitsDB.Upsert(core.Habit{Date: day.Date, State: core.HabitStateNoInfo, Type: core.HabitTypeEatHealthy})
		err = errors.Join(e1, e2, e3, e4)
		if err != nil {
			return
		}
		return daysDB.Get(core.DayParams{ID: day.ID})
	}
	return
}
