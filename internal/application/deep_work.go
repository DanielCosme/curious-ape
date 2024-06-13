package application

import (
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database"
)

func (a *App) deepWorkSync(d core.Date) error {
	dwLog, err := a.deepWorkFromToggl(d)
	if err != nil {
		return err
	}
	if !dwLog.IsZero() {
		habitLogParams := dwLog.ToHabitLogDeepWork()
		_, err := a.HabitUpsert(habitLogParams)
		if err != nil {
			return err
		}
		dwLog, err = a.db.DeepWork.Upsert(dwLog)
		if err != nil {
			return err
		}
		a.Log.Info("Deep Work log added", "date", dwLog.Date, "duration", dwLog.Duration.String())
	}
	return nil
}

func (a *App) deepWorkFromToggl(d core.Date) (res core.DeepWorkLog, err error) {
	day, err := a.db.Days.GetOrCreate(database.DayParams{Date: d})
	if err != nil {
		return res, err
	}
	summary, err := a.sync.TogglAPI.Reports.GetDaySummary(day.Date.Time())
	if err != nil {
		return
	}
	res = core.NewDeepWorkLog(summary.TotalDuration, day)
	res.IsAutomated = true
	res.Origin = core.IntegrationToggl
	res.DayID = day.ID
	return
}
