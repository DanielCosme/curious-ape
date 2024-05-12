package application

import (
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database"
	"time"
)

func (a *App) DayGetByID(id int32) (core.Day, error) {
	return a.db.Days.Get(database.DayParams{ID: id, R: database.DayRelations()})
}

// DaysCurMonth will return all the Days of the current Month.
func (a *App) DaysCurMonth() ([]core.Day, error) {
	today := core.NewDate(time.Now())
	day, err := a.db.Days.Get(database.DayParams{Date: today})
	if database.IfNotFoundErr(err) {
		return nil, err
	}

	daysOfTheMonth := today.RangeMonth()
	if day.IsZero() {
		var res []core.Day
		for _, date := range daysOfTheMonth {
			d, err := a.db.Days.GetOrCreate(database.DayParams{Date: date, R: database.DayRelations()})
			if err != nil {
				return nil, err
			}
			res = append(res, d)
		}
		return res, nil
	}

	return a.db.Days.Find(database.DayParams{Dates: daysOfTheMonth, R: database.DayRelations()})
}
