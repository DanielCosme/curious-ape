package day

import (
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/event"
	"github.com/stephenafamo/bob"
)

type Service struct {
	db  bob.DB
	bus event.Broker
}

func NewService(db bob.DB, bus event.Broker) *Service {
	return &Service{db: db, bus: bus}
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
			if _, err := s.GetOrCreate(date); err != nil {
				return nil, err
			}
		}
	}

	return find(s.db, core.DayParams{Dates: daysOfTheMonth, Order: order})
}

func (s *Service) GetOrCreate(date core.Date) (d core.Day, err error) {
	d, err = get(s.db, core.DayParams{Date: date})
	if core.IfErrNNotFound(err) {
		return
	}
	if d.IsZero() {
		d, err = new(s.db, date)
		if err != nil {
			return
		}
		s.bus.Publish(event.Event{
			Topic: event.DayCreated,
			Date:  &date,
		})
		// Reload so habits created by DayCreated subscribers are included.
		return get(s.db, core.DayParams{Date: date})
	}
	return
}
