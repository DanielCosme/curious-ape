package day

import (
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"github.com/stephenafamo/bob"
)

type Service struct {
	db bob.DB
}

func NewService(db bob.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Month(date core.Date, order core.OrderParam) ([]core.Day, error) {
	if time.Now().Month() != date.Time().Month() {
		date = date.LastDayOfTheMonth()
	}

	d, err := get(s.db, core.DayParams{Date: date})
	if core.IfErrNNotFound(err) {
		return nil, err
	}

	daysOfTheMonth := date.RangeMonth()
	if d.IsZero() {
		for _, date := range daysOfTheMonth {
			if _, err := getOrCreate(s.db, core.DayParams{Date: date}); err != nil {
				return nil, err
			}
			// TODO: Implement habit creation after this. Via Events.
		}
	}

	return find(s.db, core.DayParams{Dates: daysOfTheMonth, Order: order})
}

/*
  Next steps:
  - Create the habit Service
  	1. Make sure that Habits are initialized once a Day is created
  - Implement Navigation in Base Layout
  - Load Habits Page
*/
