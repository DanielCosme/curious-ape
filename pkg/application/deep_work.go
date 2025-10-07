package application

import (
	"context"
	"encoding/json"
	"fmt"
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
	entries, err := a.sync.TogglAPI.TimeEntries.GetDayEntries(d.Time())
	if err != nil {
		return err
	}

	var totalDuration time.Duration
	for _, entry := range entries {
		raw, err := json.Marshal(entry)
		if err != nil {
			return err
		}
		if entry.Stop.Before(d.Time()) {
			logger.Trace("skipping Toggl entry because it has not stopped")
			continue
		}
		params := core.DeepWorkLog{
			Date: day.Date,
			TimelineLog: core.TimelineLog{
				Title:     entry.Description,
				StartTime: entry.Start,
				EndTime:   entry.Stop,
			},
			Origin: core.LogOriginToggl,
			Raw:    raw,
		}
		log, err := a.db.DeepWork.Upsert(params)
		if err != nil {
			return err
		}
		duration := log.EndTime.Sub(log.StartTime)
		t := fmt.Sprintf("%s-%s (%s)", log.StartTime.Format(core.Time), log.EndTime.Format(core.Time), duration)
		logger.Info("Deep work log created: " + t)
		totalDuration += duration
	}

	habitState := core.HabitStateNotDone
	if totalDuration > time.Hour*5 {
		habitState = core.HabitStateDone
	}

	_, err = a.HabitUpsert(ctx, core.Habit{
		Date:      day.Date,
		Type:      core.HabitTypeDeepWork,
		State:     habitState,
		Note:      core.DurationToString(totalDuration),
		Automated: true})
	if err != nil {
		return err
	}
	return nil
}
