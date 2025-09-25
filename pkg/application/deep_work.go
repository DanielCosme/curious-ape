package application

import (
	"github.com/aarondl/opt/omit"
	"github.com/danielcosme/curious-ape/database/gen/models"
	"github.com/danielcosme/curious-ape/pkg/core"
	"time"
)

func (a *App) deepWorkSync(d core.Date) error {
	day, err := a.db.Days.GetOrCreate(core.DayParams{Date: d})
	if err != nil {
		return err
	}
	summary, err := a.sync.TogglAPI.Reports.GetDaySummary(day.Date.Time())
	if err != nil {
		return err
	}

	workLog, err := a.db.DeepWork.Upsert(&models.DeepWorkLogSetter{
		DayID:   omit.From(int64(day.ID)),
		Date:    omit.From(day.Date.Time()),
		Seconds: omit.From(int64(summary.TotalDuration.Seconds())),
		Origin:  omit.From(core.OriginLogToggl)})
	if err != nil {
		return err
	}
	// TODO: make this better.
	dur := time.Duration(workLog.Seconds) * time.Second
	a.Log.Info("Deep Work log added", "date", workLog.Date, "duration", dur.String())
	habitState := core.HabitStateNotDone
	if dur > time.Hour*5 {
		habitState = core.HabitStateDone
	}

	_, err = a.HabitUpsert(core.UpsertHabitParams{
		Date:      day.Date,
		Type:      core.HabitTypeDeepWork,
		State:     habitState,
		Automated: true})
	if err != nil {
		return err
	}
	return nil
}
