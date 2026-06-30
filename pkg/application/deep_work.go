package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/day"
	"danicos.dev/daniel/curious-ape/pkg/oak"
)

func (a *App) deepWorkSync(ctx context.Context, date core.Date) error {
	logger := oak.FromContext(ctx)

	if a.sync.TogglAPI == nil {
		return errors.New("Toggl API struct is nil")
	}

	d, err := day.GetOrCreate(date)
	if err != nil {
		return err
	}
	entries, err := a.sync.TogglAPI.TimeEntries.GetDayEntries(date.Time())
	if err != nil {
		return err
	}

	logger.Info("Deep work logs for: "+d.Date.Time().Format(core.HumanDateWeekDay), "entries", len(entries))
	var totalDuration time.Duration
	for _, entry := range entries {
		if entry.Stop.Before(date.ToBeginningOfDay()) {
			logger.Info("skipping Toggl entry because it has not stopped")
			continue
		}
		if entry.Start.Before(date.ToBeginningOfDay()) {
			logger.Info("skipping Toggl entry because it started before the desired day")
			continue
		}
		if entry.Start.After(date.ToEndOfDay()) {
			logger.Info("skipping Toggl entry because it is beyond the current day")
			continue
		}
		raw, err := json.Marshal(entry)
		if err != nil {
			return err
		}
		params := core.DeepWorkLog{
			Date: d.Date,
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
		logger.Info("Deep work log added: " + t)
		totalDuration += duration
	}

	habitState := core.HabitStateNotDone
	if totalDuration > time.Hour*5 {
		habitState = core.HabitStateDone
	}

	_, err = a.HabitUpsert(ctx, core.Habit{
		Date:      d.Date,
		Type:      core.HabitTypeDeepWork,
		State:     habitState,
		Note:      core.DurationToString(totalDuration),
		Automated: true})
	if err != nil {
		return err
	}
	return nil
}
