package application

import (
	"context"
	"fmt"

	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/oak"
)

func (a *App) DayGetByID(id uint) (core.Day, error) {
	return a.db.Days.Get(core.DayParams{ID: id})
}

func (a *App) DayGetOrCreate(date core.Date) (core.Day, error) {
	return a.db.Days.GetOrCreate(core.DayParams{Date: date})
}

// DaysMonth will return all the Days of the current Month.
func (a *App) DaysMonth(ctx context.Context, today core.Date) ([]core.Day, error) {
	logger := oak.FromContext(ctx, "app")
	defer logger.PopLayer()

	day, err := a.db.Days.Get(core.DayParams{Date: today})
	if core.IfErrNNotFound(err) {
		return nil, err
	}

	daysOfTheMonth := today.RangeMonth()
	logger.Info(fmt.Sprintf("Number of days: %d", len(daysOfTheMonth)))
	if day.IsZero() {
		var res []core.Day
		for _, date := range daysOfTheMonth {
			d, err := a.db.Days.GetOrCreate(core.DayParams{Date: date})
			if err != nil {
				return nil, err
			}
			res = append(res, d)
		}
		return res, nil
	}

	return a.db.Days.Find(core.DayParams{Dates: daysOfTheMonth})
}
