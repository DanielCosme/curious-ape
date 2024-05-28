package application

import (
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database"
)

func (a *App) DayGetByID(id int32) (core.Day, error) {
	return a.db.Days.Get(database.DayParams{ID: id, R: database.DayRelations()})
}

func (a *App) DayGetOrCreate(date core.Date) (core.Day, error) {
	return a.db.Days.GetOrCreate(database.DayParams{Date: date, R: database.DayRelations()})
}

// DaysMonth will return all the Days of the current Month.
func (a *App) DaysMonth(today core.Date) ([]core.Day, error) {
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

type DaysSlice []core.Day

func (a DaysSlice) Len() int           { return len(a) }
func (a DaysSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DaysSlice) Less(i, j int) bool { return a[i].Date.Time().After(a[j].Date.Time()) }
