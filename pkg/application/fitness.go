package application

import (
	"encoding/json"
	"github.com/aarondl/opt/omit"
	"github.com/danielcosme/curious-ape/database/gen/models"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/integrations/google"
)

func (a *App) fitnessSync(d core.Date) error {
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
		a.Log.Info("Fitness log added", "date", fl.Date, "origin", fl.Origin)

		if idx == 0 {
			habitState = core.HabitStateDone
		}
	}
	_, err = a.HabitUpsertAutomated(d, core.HabitTypeFitness, habitState)
	return err
}

func (a *App) fitnessLogsFromGoogle(d core.Date) (res []*models.FitnessLogSetter, err error) {
	googleClient, err := a.googleClient()
	if err != nil {
		return
	}
	day, err := a.DayGetOrCreate(d)
	if err != nil {
		return
	}

	date := core.NewDate(day.Date)
	sessions, err := googleClient.Fitness.GetFitnessSessions(date.ToBeginningOfDay(), date.ToEndOfDay())
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

func fitnessLogFromGoogle(day *models.Day, session google.Session) (*models.FitnessLogSetter, error) {
	raw, err := json.Marshal(&session)
	if err != nil {
		return nil, err
	}
	setter := &models.FitnessLogSetter{
		DayID:     omit.From(day.ID),
		Type:      omit.From("strong"),
		Title:     omit.From(session.Name),
		Date:      omit.From(day.Date),
		StartTime: omit.From(google.ParseMillis(session.StartTimeMillis)),
		EndTime:   omit.From(google.ParseMillis(session.EndTimeMillis)),
		Origin:    omit.From(core.OriginLogGoogle),
		Note:      omit.From(session.Application.PackageName),
		Raw:       omit.From(string(raw)),
	}
	return setter, nil
}
