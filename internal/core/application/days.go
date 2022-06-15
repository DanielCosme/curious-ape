package application

import (
	"github.com/danielcosme/curious-ape/internal/core/database"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/sdk/errors"
	"time"
)

func (a *App) DayCreate(d *entity.Day) (*entity.Day, error) {
	d.Date = time.Date(d.Date.Year(), d.Date.Month(), d.Date.Day(), 0, 0, 0, 0, time.UTC)
	if err := a.db.Days.Create(d); err != nil {
		return nil, err
	}

	return a.DayGetByDate(d.Date)
}

func (a *App) DaysGetAll() ([]*entity.Day, error) {
	return a.db.Days.Find(entity.DayFilter{}, database.DaysPipeline(a.db)...)
}

func (a *App) DayGetByDate(date time.Time) (*entity.Day, error) {
	d, err := a.db.Days.Get(entity.DayFilter{Dates: []time.Time{date}})
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return nil, err
	}
	if d == nil {
		// if it does not exist, create new and return.
		return a.DayCreate(&entity.Day{Date: date})
	}

	return d, nil
}

func (a *App) daysGetByDateRange(start, end time.Time) ([]*entity.Day, error) {
	days := []*entity.Day{}

	for _, date := range datesRange(start, end) {
		d, err := a.DayGetByDate(date)
		if err != nil {
			return nil, err
		}

		days = append(days, d)
	}

	return days, nil
}

func datesRange(start, end time.Time) []time.Time {
	dates := []time.Time{}
	inter := start

	for inter.Before(end) {
		dates = append(dates, inter)
		inter = inter.AddDate(0, 0, 1)
	}
	dates = append(dates, end)

	return dates
}
