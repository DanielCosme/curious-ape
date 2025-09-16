package application

import (
	"github.com/danielcosme/curious-ape/database/gen/models"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/persistence"
)

func (a *App) DayGetByID(id int64) (*models.Day, error) {
	return a.db.Days.Get(persistence.DayParams{ID: id, R: persistence.DayRelations()})
}

func (a *App) DayGetOrCreate(date core.Date) (*models.Day, error) {
	return a.db.Days.GetOrCreate(persistence.DayParams{Date: date, R: persistence.DayRelations()})
}

// DaysMonth will return all the Days of the current Month.
func (a *App) DaysMonth(today core.Date) ([]*models.Day, error) {
	day, err := a.db.Days.Get(persistence.DayParams{Date: today})
	if persistence.IgnoreIfErrNotFound(err) {
		return nil, err
	}

	daysOfTheMonth := today.RangeMonth()
	if day == nil {
		var res []*models.Day
		for _, date := range daysOfTheMonth {
			d, err := a.db.Days.GetOrCreate(persistence.DayParams{Date: date, R: persistence.DayRelations()})
			if err != nil {
				return nil, err
			}
			res = append(res, d)
		}
		return res, nil
	}

	return a.db.Days.Find(persistence.DayParams{Dates: daysOfTheMonth, R: persistence.DayRelations()})
}
