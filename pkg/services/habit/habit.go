package habit

import (
	"fmt"
	"log/slog"
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/event"
	"github.com/nats-io/nats.go"
	"github.com/stephenafamo/bob"
)

type Service struct {
	db   bob.DB
	nats *nats.Conn
}

func NewService(db bob.DB, ns *nats.Conn) *Service {
	s := &Service{
		db:   db,
		nats: ns,
	}
	if ns != nil {
		ns.Subscribe(event.DayCreated, s.listen)
		ns.Subscribe(event.WorklogSynced, s.listen)
		ns.Subscribe(event.FitnesslogSynced, s.listen)
		ns.Subscribe(event.SleepLogSynced, s.listen)
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
		var wls core.LogSyncPayload
		core.Decode(msg.Data, &wls)
		if err := s.eventWorklogSynced(wls); err != nil {
			slog.Error("Failed to handle event", "err", err, "subject", msg.Subject)
		}
	case event.FitnesslogSynced:
		var wls core.LogSyncPayload
		core.Decode(msg.Data, &wls)
		if err := s.eventFitnesslogSynced(wls); err != nil {
			slog.Error("Failed to handle event", "err", err, "subject", msg.Subject)
		}
	case event.SleepLogSynced:
		var wls core.LogSyncPayload
		core.Decode(msg.Data, &wls)
		if err := s.eventSleepLogSynced(wls); err != nil {
			slog.Error("Failed to handle event", "err", err, "subject", msg.Subject)
		}
	}
}

func (svc *Service) eventWorklogSynced(payload core.LogSyncPayload) error {
	habitParams := core.Habit{
		Date:      payload.Date,
		Type:      core.HabitTypeDeepWork,
		State:     core.HabitStateNotDone,
		Automated: true,
	}

	var totalDuration time.Duration
	for _, wl := range payload.WorkLogs {
		duration := wl.EndTime.Sub(wl.StartTime)
		totalDuration += duration
	}

	if totalDuration > time.Hour*5 {
		habitParams.State = core.HabitStateDone
	}
	if totalDuration > 0 {
		habitParams.Note = core.DurationToString(totalDuration)
	}

	_, err := svc.HabitUpsert(habitParams)
	if err != nil {
		return err
	}
	svc.nats.Publish(event.DaySynced, payload.Date.Enc())
	return nil
}

func (svc *Service) eventFitnesslogSynced(payload core.LogSyncPayload) error {
	habitParams := core.Habit{
		Date:      payload.Date,
		Type:      core.HabitTypeFitness,
		State:     core.HabitStateNotDone,
		Automated: true,
	}

	for idx, fl := range payload.FitnessLogs {
		if idx == 0 {
			habitParams.State = core.HabitStateDone
			duration := core.DurationToString(fl.EndTime.Sub(fl.StartTime))
			habitParams.Note = fmt.Sprintf("%s - %s (%s)", fl.StartTime.Format(core.Time), fl.EndTime.Format(core.Time), duration)
			break
		}
	}

	_, err := svc.HabitUpsert(habitParams)
	if err != nil {
		return err
	}
	svc.nats.Publish(event.DaySynced, payload.Date.Enc())
	return nil
}

func (svc *Service) eventSleepLogSynced(payload core.LogSyncPayload) error {
	for _, sl := range payload.SleepLogs {
		if sl.IsMainSleep {
			habitState := core.HabitStateNotDone
			wakeUpTime := time.Date(sl.EndTime.Year(), sl.EndTime.Month(), sl.EndTime.Day(), 6, 0, 0, 0, sl.EndTime.Location())
			if sl.EndTime.Before(wakeUpTime) {
				habitState = core.HabitStateDone
			}
			params := core.Habit{
				Date:      payload.Date,
				Type:      core.HabitTypeWakeUp,
				State:     habitState,
				Note:      sl.EndTime.Format(core.Time),
				Automated: true,
			}
			_, err := svc.HabitUpsert(params)
			if err != nil {
				return err
			}
		}
	}

	svc.nats.Publish(event.DaySynced, payload.Date.Enc())
	return nil
}
