package day

import (
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/event"
	"danicos.dev/daniel/curious-ape/pkg/utils"
	"github.com/nats-io/nats.go"
	"github.com/stephenafamo/bob"
)

type Service struct {
	db   bob.DB
	nats *nats.Conn
}

func NewService(db bob.DB, ns *nats.Conn) *Service {
	return &Service{db: db, nats: ns}
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

	return Find(s.db, core.DayParams{
		Dates:        daysOfTheMonth,
		Order:        order,
		WithRelation: []core.DayRelations{core.DayRelationHabits},
	},
	)
}

func (s *Service) GetOrCreate(date core.Date) (d core.Day, err error) {
	params := core.DayParams{Date: date, WithRelation: []core.DayRelations{core.DayRelationHabits}}
	d, err = get(s.db, params)
	if core.IfErrNNotFound(err) {
		return
	}
	if d.IsZero() {
		d, err = new(s.db, date)
		if err != nil {
			return
		}

		// NOTE: the first subscriber to respond will un-block this. Might need to address this in the future.
		_, err := s.nats.Request(event.DayCreated, date.Enc(), time.Second*10)
		if err != nil && !utils.ErrNatsIsNoResponders(err) {
			return d, err
		}

		return get(s.db, params)
	}
	return
}
