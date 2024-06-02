package application

import "github.com/danielcosme/curious-ape/internal/core"

func (a *App) DeepWorkSync(d core.Date) error {
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
		// TODO: Insert Deep Work Log into the database, also log it.
	}
	return nil
}

func (a *App) deepWorkFromToggl(d core.Date) (res core.DeepWorkLog, err error) {
	summary, err := a.sync.TogglAPI.Reports.GetDaySummary(d.Time())
	if err != nil {
		return
	}
	return core.NewDeepWorkLog(summary.TotalDuration, d), err
}
