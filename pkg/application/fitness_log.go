package application

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/integrations/google"
	"github.com/danielcosme/curious-ape/pkg/oak"
)

func (a *App) fitnessSync(ctx context.Context, d core.Date) error {
	logger := oak.FromContext(ctx)

	fitnessLogs, err := a.fitnessLogsFromHevy(ctx, d)
	if err != nil {
		return err
	}

	habitParams := core.Habit{
		Date:      d,
		Type:      core.HabitTypeFitness,
		State:     core.HabitStateNotDone,
		Automated: true,
	}
	for idx, fl := range fitnessLogs {
		fl, err := a.db.Fitness.Upsert(fl)
		if err != nil {
			return err
		}
		logger.Info("Fitness log added", "date", fl.Date, "origin", fl.Origin)

		if idx == 0 {
			habitParams.State = core.HabitStateDone

			duration := core.DurationToString(fl.EndTime.Sub(fl.StartTime))
			habitParams.Note = fmt.Sprintf("%s - %s (%s)", fl.StartTime.Format(core.Time), fl.EndTime.Format(core.Time), duration)
		}
	}
	_, err = a.HabitUpsert(ctx, habitParams)
	return err
}

func (a *App) fitnessLogsFromHevy(ctx context.Context, d core.Date) (res []core.FitnessLog, err error) {
	logger := oak.FromContext(ctx)

	if d.Time().Before(core.NewDate(time.Now()).Time()) {
		//NOTE: no-op if the desired log is not for the current day. TODO: support this in the future (their API is funny).
		logger.Warning("fitness log to sync is not today", "date", d.String())
		return res, fmt.Errorf("fitness log to sync is not today: %s", d.String())
	}

	day, err := a.dayGetOrCreate(d)
	if err != nil {
		return nil, err
	}
	events, err := a.sync.Hevy.WorkoutEvents.Get(d.Time())
	if err != nil {
		return nil, err
	}

	logger.Info("Fitness log for: "+day.Date.Time().Format(core.HumanDateWeekDay), "entries", len(events))
	for _, event := range events {
		if event.Type != "updated" {
			return res, fmt.Errorf("unkouwn event type: %s", event.Type)
		}

		raw, err := json.Marshal(event.Workout)
		if err != nil {
			return res, err
		}
		fitnessLogType := core.FitnessLogTypeOther
		title := strings.ToLower(event.Workout.Title)
		if strings.Contains(title, "lower") || strings.Contains(title, "upper") {
			fitnessLogType = core.FitnessLogTypeStrength
		} else if strings.Contains(title, "cardio") {
			fitnessLogType = core.FitnessLogTypeCardio
		}

		normalizeTime := func(t time.Time, loc *time.Location) time.Time {
			return core.TimeUTC(t.In(loc))
		}
		location, _ := time.LoadLocation("America/Toronto")
		fl := core.FitnessLog{
			Date: day.Date,
			TimelineLog: core.TimelineLog{
				Title:     event.Workout.Title,
				StartTime: normalizeTime(event.Workout.StartTime, location),
				EndTime:   normalizeTime(event.Workout.EndTime, location),
			},
			FitnessType: fitnessLogType,
			Origin:      core.LogOriginHevy,
			Raw:         raw,
		}
		res = append(res, fl)
	}
	return res, nil
}

func (a *App) fitnessLogsFromGoogle(d core.Date) (res []core.FitnessLog, err error) {
	googleClient, err := a.googleClient()
	if err != nil {
		return
	}
	day, err := a.DayGetOrCreate(d)
	if err != nil {
		return
	}

	sessions, err := googleClient.Fitness.GetFitnessSessions(day.Date.ToBeginningOfDay(), day.Date.ToEndOfDay())
	if err != nil {
		return
	}
	for _, s := range sessions {
		fl, err := fitnessLogFromGoogle(day, s)
		if err != nil {
			return nil, err
		}
		res = append(res, fl)
	}
	return
}

func fitnessLogFromGoogle(day core.Day, session google.FitnessSession) (fl core.FitnessLog, err error) {
	raw, err := json.Marshal(&session)
	if err != nil {
		return
	}

	fl.Date = day.Date
	fl.Title = session.Name
	fl.StartTime = google.ParseMillis(session.StartTimeMillis)
	fl.EndTime = google.ParseMillis(session.EndTimeMillis)
	fl.Note = session.Application.PackageName
	// TODO: find a way to discern type, right now is hardcoded to Strength
	fl.FitnessType = core.FitnessLogTypeStrength
	fl.Origin = core.LogOriginGoogle
	fl.Raw = raw
	return
}
