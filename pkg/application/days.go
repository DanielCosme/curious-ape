package application

import (
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/database"
	"github.com/danielcosme/curious-ape/pkg/database/gen/models"
)

func (a *App) DayGetByID(id int32) (*models.Day, error) {
	return a.db.Days.Get(database.DayParams{ID: id, R: database.DayRelations()})
}

func (a *App) DayGetOrCreate(date core.Date) (*models.Day, error) {
	return a.db.Days.GetOrCreate(database.DayParams{Date: date, R: database.DayRelations()})
}

// DaysMonth will return all the Days of the current Month.
func (a *App) DaysMonth(today core.Date) ([]*models.Day, error) {
	day, err := a.db.Days.Get(database.DayParams{Date: today})
	if database.IgnoreIfErrNotFound(err) {
		return nil, err
	}

	daysOfTheMonth := today.RangeMonth()
	if day == nil {
		var res []*models.Day
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
