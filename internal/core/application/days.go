package application

import (
	"errors"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/datasource"
	"github.com/danielcosme/curious-ape/sdk/dates"
	"time"
)

func (a *App) DayCreate(d *entity.Day) (*entity.Day, error) {
	var err error
	d.Date = dates.ToUTC(d.Date)
	if err = a.db.Days.Create(d); err != nil {
		return nil, err
	}

	return a.DayGetByDate(d.Date)
}

func (a *App) DaysGetAll() ([]*entity.Day, error) {
	return a.db.Days.Find(entity.DayFilter{}, datasource.DaysPipeline(a.db)...)
}

func (a *App) DayGetByDate(date time.Time) (*entity.Day, error) {
	d, err := a.db.Days.Get(entity.DayFilter{Date: []time.Time{date}})
	if err != nil && errors.Is(err, entity.ErrNotFound) {
		// if it does not exist, create new and return.
		return a.DayCreate(&entity.Day{Date: date})
	}

	return d, err
}
