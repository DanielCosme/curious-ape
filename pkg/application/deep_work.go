package application

import (
	"context"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/oak"
)

func (a *App) deepWorkSync(ctx context.Context, d core.Date) error {
	logger := oak.FromContext(ctx)

	day, err := a.dayGetOrCreate(d)
	if err != nil {
		return err
	}
	summary, err := a.sync.TogglAPI.Reports.GetDaySummary(day.Date.Time())
	if err != nil {
		return err
	}
	logger.Notice("Deep work sync not implemented", "Work-log-total-duration", summary.TotalDuration.String())

	/*
		workLog, err := a.db.DeepWork.Upsert(&models.DeepWorkLogSetter{
			Title:     omit.From("Deep Work"),
			DayID:     omit.From(int64(day.ID)),
			Date:      omit.From(day.Date.Time()),
			Seconds:   omit.From(int64(summary.TotalDuration.Seconds())),
			StartTime: omit.From(time.Now()),
			Raw:       omit.From(""),
			Origin:    omit.From(core.OriginLogToggl)})
		if err != nil {
			return err
		}
		// TODO: make this better.
		dur := time.Duration(workLog.Seconds) * time.Second
		logger.Info("Deep Work log added", "date", workLog.Date, "duration", dur.String())
		habitState := core.HabitStateNotDone
		if dur > time.Hour*5 {
			habitState = core.HabitStateDone
		}

		_, err = a.HabitUpsert(ctx, core.UpsertHabitParams{
			Date:      day.Date,
			Type:      core.HabitTypeDeepWork,
			State:     habitState,
			Automated: true})
		if err != nil {
			return err
		}
	*/
	return nil
}
