package habit

import (
	"log/slog"
	"time"

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
		ns.Subscribe(event.WorklogSynced, s.listen)
	}
	return s
}

func (s *Service) HabitUpsert(params core.Habit) (habit core.Habit, err error) {
	habit, err = upsert(s.db, params)
	if err == nil && habit.State != core.HabitStateNoInfo {
		slog.Info("Habit logged",
			"type", habit.Type,
			"state", habit.State,
			"date", habit.Date.Time.Format(core.HumanDateWeekDay),
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
				slog.Error("Failed to create habit", "err", err)
			}
		}
		msg.Respond(nil)
	case event.WorklogSynced:
		var wls []core.DeepWorkLog
		core.Decode(msg.Data, &wls)
		if err := s.eventWorklogSynced(wls); err != nil {
			slog.Error("Failed to handle event", "err", err, "subject", msg.Subject)
		}
	}
}

func (s *Service) eventWorklogSynced(wls []core.DeepWorkLog) error {
	var date core.Date
	if len(wls) > 0 {
		date = wls[0].Date
	} else {
		slog.Info("habit: no work logs to handle")
		return nil
	}

	var totalDuration time.Duration
	for _, wl := range wls {
		duration := wl.EndTime.Sub(wl.StartTime)
		totalDuration += duration
	}
	habitState := core.HabitStateNotDone
	if totalDuration > time.Hour*5 {
		habitState = core.HabitStateDone
	}
	habitParams := core.Habit{
		Date:      date,
		Type:      core.HabitTypeDeepWork,
		State:     habitState,
		Note:      core.DurationToString(totalDuration),
		Automated: true,
	}
	_, err := s.HabitUpsert(habitParams)
	if err != nil {
		return err
	}
	s.ns.Publish(event.DaySynced, date.Enc())
	return nil
}
