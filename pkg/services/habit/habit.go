package habit

import (
	"log/slog"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/event"
	"github.com/nats-io/nats.go"
	"github.com/stephenafamo/bob"
)

type Service struct {
	db bob.DB
	ns *nats.Conn
}

func NewService(db bob.DB, ns *nats.Conn) *Service {
	s := &Service{
		db: db,
		ns: ns,
	}
	if ns != nil {
		ns.Subscribe(event.DayCreated, s.listen)
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

func (s *Service) Flip(id int) (habit core.Habit, err error) {
	habit, err = get(s.db, core.HabitParams{ID: id})
	if err != nil {
		return
	}
	state := core.HabitStateNotDone
	if habit.State == core.HabitStateNotDone || habit.State == core.HabitStateNoInfo {
		state = core.HabitStateDone
	}
	habit.State = state
	return upsert(s.db, core.Habit{
		Date:  habit.Date,
		Type:  habit.Type,
		State: habit.State,
	})
}

func (s *Service) listen(msg *nats.Msg) {
	switch msg.Subject {
	case event.DayCreated:
		date := core.DateDecode(msg.Data)
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
		msg.Respond(nil)
	}
}
