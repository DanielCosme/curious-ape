package application

import (
	"context"
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

	habitParams := core.UpsertHabitParams{
		Date:      d,
		Type:      core.HabitTypeFitness,
		State:     core.HabitStateNotDone,
		Automated: true,
	}
	for idx, fl := range fitnessLogs {
		// fl, err := a.db.Fitness.Upsert(setter)
		// if err != nil {
		// 	return err
		// }
		// logger.Info("Fitness log added", "date", fl.Date, "origin", fl.Origin)

		if idx == 0 {
			habitParams.State = core.HabitStateDone

			logger.Notice("no fitness log is going to be created")
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
	// raw, err := json.Marshal(&session)
	// if err != nil {
	// 	return
	// }

	fl.StartTime = google.ParseMillis(session.StartTimeMillis)
	fl.EndTime = google.ParseMillis(session.EndTimeMillis)

	// setter := &models.FitnessLogSetter{
	// 	DayID:     omit.From(int64(day.ID)),
	// 	Type:      omit.From("strong"),
	// 	Title:     omit.From(session.Name),
	// 	Date:      omit.From(day.Date.Time()),
	// 	StartTime: omit.From(google.ParseMillis(session.StartTimeMillis)),
	// 	EndTime:   omit.From(google.ParseMillis(session.EndTimeMillis)),
	// 	Origin:    omit.From(core.OriginLogGoogle),
	// 	Note:      omit.From(session.Application.PackageName),
	// 	Raw:       omit.From(string(raw)),
	// }
	return
}
