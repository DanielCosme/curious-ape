package application

import (
	"context"
	"time"

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

	logger.Notice("no deep work log created")
	// workLog, err := a.db.DeepWork.Upsert(&models.DeepWorkLogSetter{
	// 	Title:     omit.From("Deep Work"),
	// 	DayID:     omit.From(int64(day.ID)),
	// 	Date:      omit.From(day.Date.Time()),
	// 	Seconds:   omit.From(int64(summary.TotalDuration.Seconds())),
	// 	StartTime: omit.From(time.Now()),
	// 	Raw:       omit.From(""),
	// 	Origin:    omit.From(core.OriginLogToggl)})
	// if err != nil {
	// 	return err
	// }

	habitState := core.HabitStateNotDone
	if summary.TotalDuration > time.Hour*5 {
		habitState = core.HabitStateDone
	}

	_, err = a.HabitUpsert(ctx, core.Habit{
		Date:      day.Date,
		Type:      core.HabitTypeDeepWork,
		State:     habitState,
		Note:      core.DurationToString(summary.TotalDuration),
		Automated: true})
	if err != nil {
		return err
	}
	return nil
}
