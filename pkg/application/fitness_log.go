package application

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/integrations/google"
	"github.com/danielcosme/curious-ape/pkg/oak"
)

func (a *App) fitnessSync(ctx context.Context, d core.Date) error {
	logger := oak.FromContext(ctx)

	fitnessLogs, err := a.fitnessLogsFromGoogle(d)
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
