package application

import (
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/database"
	"github.com/danielcosme/curious-ape/pkg/database/gen/models"
	"time"
)

func (a *App) deepWorkSync(d core.Date) error {
	day, err := a.db.Days.GetOrCreate(database.DayParams{Date: d})
	if err != nil {
		return err
	}
	summary, err := a.sync.TogglAPI.Reports.GetDaySummary(day.Date)
	if err != nil {
		return err
	}

	workLog, err := a.db.DeepWork.Upsert(&models.DeepWorkLogSetter{
		DayID:       omit.From(day.ID),
		Date:        omit.From(day.Date),
		Seconds:     omit.From(int32(summary.TotalDuration.Seconds())),
		IsAutomated: omitnull.From(true),
		Origin:      omit.From(core.OriginLogToggl),
	})
	if err != nil {
		return err
	}
	// TODO(daniel) make this better.
	dur := time.Duration(workLog.Seconds) * time.Second
	a.Log.Info("Deep Work log added", "date", workLog.Date, "duration", dur.String())
	habitState := core.HabitStateNotDone
	if dur > time.Hour*5 {
		habitState = core.HabitStateDone
	}

	_, err = a.HabitUpsert(core.NewDate(day.Date), core.HabitTypeDeepWork, habitState)
	if err != nil {
		return err
	}
	return nil
}
