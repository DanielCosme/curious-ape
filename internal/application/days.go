package application

import (
	"errors"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database"
	"time"
)

// DaysCurMonth will return all the Days of the current Month.
func (a *App) DaysCurMonth() ([]*core.Day, error) {
	var res []*core.Day

	today := core.NewDate(time.Now())
	day, err := a.db.Days.Get(database.DayF{Date: today})
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return nil, err
	}

	daysOfTheMonth := today.RangeMonth()
	if day != nil {
		res, err = a.db.Days.Find(database.DayF{Dates: daysOfTheMonth, WithAll: true})
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	for _, date := range daysOfTheMonth {
		d, err := a.db.Days.GetOrCreate(database.DayF{Date: date})
		if err != nil {
			return nil, err
		}
		res = append(res, d)
	}
	return res, nil
}
