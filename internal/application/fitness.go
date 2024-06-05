package application

import (
	"encoding/json"
	"errors"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/integrations/google"
)

func (a *App) FitnessSync(d core.Date) error {
	fitnessLogs, err := a.fitnessLogsFromGoogle(d)
	if err != nil {
		return err
	}
	for idx, fl := range fitnessLogs {
		if idx == 0 {
			_, err := a.HabitUpsert(fl.ToHabitLogFitness())
			if err != nil {
				return err
			}
		}
		fl, err = a.db.Fitness.Upsert(fl)
		if err != nil {
			return err
		}
		a.Log.Info("Fitness log added", "date", fl.Date, "origin", fl.Origin)
	}
	if len(fitnessLogs) == 0 {
		_, err := a.HabitUpsert(core.NewHabitParams{
			Success:   false,
			Date:      d,
			HabitType: core.HabitTypeExercise,
			Origin:    core.OriginLogFitness,
			Automated: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
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
		fitnessLog, err := fitnessLogFromGoogle(day, s)
		if err != nil {
			return nil, err
		}
		res = append(res, fitnessLog)
	}
	return
}

func fitnessLogFromGoogle(day core.Day, session google.Session) (res core.FitnessLog, err error) {
	startTime := google.ParseMillis(session.StartTimeMillis)
	if !day.Date.IsEqual(startTime) {
		return res, errors.New("fitness log from google: expected date-time for current day do not match")
	}

	res = core.NewFitnessLog(day)
	res.Title = session.Name
	res.Date = day.Date
	res.StartTime = startTime
	res.EndTime = google.ParseMillis(session.EndTimeMillis)
	res.Origin = core.IntegrationGoogle
	res.Note = session.Application.PackageName

	raw, err := json.Marshal(&session)
	if err != nil {
		return res, err
	}
	res.Raw = string(raw)
	return
}
