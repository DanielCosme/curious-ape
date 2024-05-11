package application

import (
	"errors"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database"
	"time"
)

// DaysCurMonth will return all the Days of the current Month.
func (a *App) DaysCurMonth() ([]core.Day, error) {
	var res []core.Day

	today := core.NewDate(time.Now())
	day, err := a.db.Days.Get(database.DayParams{Date: today})
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return nil, err
	}

	daysOfTheMonth := today.RangeMonth()
	if day.IsZero() {
		res, err = a.db.Days.Find(database.DayParams{
			Dates: daysOfTheMonth,
			R:     database.DayRelations(),
		})
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	for _, date := range daysOfTheMonth {
		d, err := a.db.Days.GetOrCreate(database.DayParams{Date: date, R: database.DayRelations()})
		if err != nil {
			return nil, err
		}
		res = append(res, d)
	}
	return res, nil
}
