package habit

import (
	"log/slog"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/event"
	"github.com/stephenafamo/bob"
)

type Service struct {
	db     bob.DB
	events event.EventChan
}

func NewService(db bob.DB, bus event.Broker) *Service {
	s := &Service{
		db:     db,
		events: make(event.EventChan),
	}
	if bus != nil {
		bus.Subscribe(event.DayCreated, s.events)
		go s.Listen()
	}
	return s
}

func (s *Service) HabitUpsert(params core.Habit) (habit core.Habit, err error) {
	habit, err = upsert(s.db, params)
	if err == nil && habit.State != core.HabitStateNoInfo {
		slog.Info("Habit logged",
			"type", habit.Type,
			"state", habit.State,
			"date", habit.Date.Time().Format(core.HumanDateWeekDay),
		)
	}
	return
}

func (s *Service) Listen() {
	for {
		ev, ok := <-s.events
		if !ok {
			return
		}

		switch ev.Topic {
		case event.DayCreated:
			if ev.Date == nil {
				return
			}
			date := *ev.Date
			params := []core.Habit{
				{Date: date, State: core.HabitStateNoInfo, Type: core.HabitTypeWakeUp},
				{Date: date, State: core.HabitStateNoInfo, Type: core.HabitTypeFitness},
				{Date: date, State: core.HabitStateNoInfo, Type: core.HabitTypeDeepWork},
				{Date: date, State: core.HabitStateNoInfo, Type: core.HabitTypeEatHealthy},
			}
			for _, param := range params {
				if _, err := s.HabitUpsert(param); err != nil {
					slog.Error("Failed to upsert habit", "err", err)
				}
			}

			// Unblock Publish only after work for this event is done.
			ev.Done()
		}
	}
}
