package application

import (
	"context"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/oak"
)

func (a *App) fitnessSync(ctx context.Context, d core.Date) error {
	logger := oak.FromContext(ctx)
	logger.Notice("Fitness logging not implemented")

	/*
		fitnessLogs, err := a.fitnessLogsFromGoogle(d)
		if err != nil {
			return err
		}

		habitState := core.HabitStateNotDone
		for idx, setter := range fitnessLogs {
			fl, err := a.db.Fitness.Upsert(setter)
			if err != nil {
				return err
			}
			logger.Info("Fitness log added", "date", fl.Date, "origin", fl.Origin)

			if idx == 0 {
				habitState = core.HabitStateDone
			}
		}
		_, err = a.HabitUpsert(ctx, core.UpsertHabitParams{
			Date:      d,
			Type:      core.HabitTypeFitness,
			State:     habitState,
			Automated: true})
		return err
	*/
	return nil
}

/*
func (a *App) fitnessLogsFromGoogle(d core.Date) (res []*models.FitnessLogSetter, err error) {
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

func fitnessLogFromGoogle(day core.Day, session google.Session) (*models.FitnessLogSetter, error) {
	raw, err := json.Marshal(&session)
	if err != nil {
		return nil, err
	}
	setter := &models.FitnessLogSetter{
		DayID:     omit.From(int64(day.ID)),
		Type:      omit.From("strong"),
		Title:     omit.From(session.Name),
		Date:      omit.From(day.Date.Time()),
		StartTime: omit.From(google.ParseMillis(session.StartTimeMillis)),
		EndTime:   omit.From(google.ParseMillis(session.EndTimeMillis)),
		Origin:    omit.From(core.OriginLogGoogle),
		Note:      omit.From(session.Application.PackageName),
		Raw:       omit.From(string(raw)),
	}
	return setter, nil
}
*/
