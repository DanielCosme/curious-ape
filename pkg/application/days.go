package application

import (
	"context"
	"errors"

	"github.com/danielcosme/curious-ape/pkg/core"
)

func (a *App) DayGetByID(id uint) (core.Day, error) {
	return a.db.Days.Get(core.DayParams{ID: id})
}

func (a *App) DayGetOrCreate(date core.Date) (core.Day, error) {
	return a.dayGetOrCreate(date)
}

// DaysMonth will return all the Days of the current Month.
func (a *App) DaysMonth(ctx context.Context, today core.Date) ([]core.Day, error) {
	day, err := a.db.Days.Get(core.DayParams{Date: today})
	if core.IfErrNNotFound(err) {
		return nil, err
	}

	daysOfTheMonth := today.RangeMonth()
	if day.IsZero() {
		var res []core.Day
		for _, date := range daysOfTheMonth {
			d, err := a.dayGetOrCreate(date)
			if err != nil {
				return nil, err
			}
			res = append(res, d)
		}
		return res, nil
	}

	return a.db.Days.Find(core.DayParams{Dates: daysOfTheMonth})
}

func (a *App) dayGetOrCreate(d core.Date) (day core.Day, err error) {
	day, err = a.db.Days.Get(core.DayParams{Date: d})
	if core.IfErrNNotFound(err) {
		return
	}
	if day.IsZero() {
		day, err = a.db.Days.Create(d)
		if err != nil {
			return
		}
		_, e1 := a.db.Habits.Upsert(core.UpsertHabitParams{Date: day.Date, State: core.HabitStateNoInfo, Type: core.HabitTypeWakeUp})
		_, e2 := a.db.Habits.Upsert(core.UpsertHabitParams{Date: day.Date, State: core.HabitStateNoInfo, Type: core.HabitTypeFitness})
		_, e3 := a.db.Habits.Upsert(core.UpsertHabitParams{Date: day.Date, State: core.HabitStateNoInfo, Type: core.HabitTypeDeepWork})
		_, e4 := a.db.Habits.Upsert(core.UpsertHabitParams{Date: day.Date, State: core.HabitStateNoInfo, Type: core.HabitTypeEatHealthy})
		err = errors.Join(e1, e2, e3, e4)
		if err != nil {
			return
		}
		return a.db.Days.Get(core.DayParams{ID: day.ID})
	}
	return
}
