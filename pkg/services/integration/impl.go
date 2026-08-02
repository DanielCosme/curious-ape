package integration

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"danicos.dev/daniel/curious-ape/pkg/config"
	"danicos.dev/daniel/curious-ape/pkg/core"
)

func (svc *Service) GetDayEntries(date core.Date) ([]core.DeepWorkLog, error) {
	if svc.sync.TogglAPI == nil {
		return nil, errors.New("Toggl API service is nil")
	}

	entries, err := svc.sync.TogglAPI.TimeEntries.GetDayEntries(date.Time)
	if err != nil {
		return nil, err
	}

	res := make([]core.DeepWorkLog, 0, len(entries))
	slog.Info("Deep work logs for: "+date.Time.Format(core.HumanDateWeekDay), "entries", len(entries))
	for _, entry := range entries {
		if entry.Stop.Before(date.ToBeginningOfDay()) {
			slog.Info("skipping Toggl entry because it has not stopped")
			continue
		}
		if entry.Start.Before(date.ToBeginningOfDay()) {
			slog.Info("skipping Toggl entry because it started before the desired day")
			continue
		}
		if entry.Start.After(date.ToEndOfDay()) {
			slog.Info("skipping Toggl entry because it is beyond the current day")
			continue
		}
		raw, err := json.Marshal(entry)
		if err != nil {
			return nil, err
		}
		workLog := core.DeepWorkLog{
			Date: date,
			TimelineLog: core.TimelineLog{
				Title:     entry.Description,
				StartTime: entry.Start,
				EndTime:   entry.Stop,
			},
			Origin: core.LogOriginToggl,
			Raw:    raw,
		}
		res = append(res, workLog)
	}

	return res, nil
}

func (svc *Service) GetFitnessLogs(date core.Date) (res []core.FitnessLog, err error) {
	if svc.sync.Hevy == nil {
		slog.Warn("Fitness provider Hevy is nil, cannot sync Fitness Logs")
		return
	}

	if !date.Time.Before(core.NewDate(time.Now()).Time) {
		events, err := svc.sync.Hevy.WorkoutEvents.Get(date.Time)
		if err == nil {
			slog.Info("Fitness log for: "+date.Time.Format(core.HumanDateWeekDay), "entries", len(events))
			for _, ev := range events {
				if ev.Type == "updated" {
					raw, err := json.Marshal(ev.Workout)
					if err == nil {
						fitnessLogType := core.FitnessLogTypeOther
						title := strings.ToLower(ev.Workout.Title)
						if strings.Contains(title, "lower") || strings.Contains(title, "upper") {
							fitnessLogType = core.FitnessLogTypeStrength
						} else if strings.Contains(title, "cardio") {
							fitnessLogType = core.FitnessLogTypeCardio
						}

						normalizeTime := func(t time.Time, loc *time.Location) time.Time {
							return core.TimeUTC(t.In(loc))
						}
						location, _ := time.LoadLocation(config.TZ)
						fl := core.FitnessLog{
							Date: date,
							TimelineLog: core.TimelineLog{
								Title:     ev.Workout.Title,
								StartTime: normalizeTime(ev.Workout.StartTime, location),
								EndTime:   normalizeTime(ev.Workout.EndTime, location),
							},
							FitnessType: fitnessLogType,
							Origin:      core.LogOriginHevy,
							Raw:         raw,
						}
						res = append(res, fl)
					} else {
						return nil, err
					}
				} else {
					return nil, fmt.Errorf("unkouwn event type: %s", ev.Type)
				}
			}
			return res, nil
		}
		return nil, err
	} else {
		slog.Warn(fmt.Sprintf("fitness log to sync is not today: %s", date.String()))
	}
	return
}
