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
	"danicos.dev/daniel/curious-ape/pkg/integrations/fitbit"
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

func (svc *Service) GetSleeLogs(date core.Date) (res []core.SleepLog, err error) {
	fitbitClient, err := svc.fitbitClient()
	if err != nil {
		return
	}
	sleepLogs, err := fitbitClient.Sleep.GetByDate(date.Time)
	if err != nil {
		return res, err
	}
	for _, fsl := range sleepLogs.Sleep {
		sl, err := sleepLogFromFitbit(date, fsl)
		if err != nil {
			return res, err
		}
		res = append(res, sl)
	}
	return
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

func sleepLogFromFitbit(date core.Date, s fitbit.Sleep) (sl core.SleepLog, err error) {
	if !date.IsEqual(fitbit.ParseDate(s.DateOfSleep)) {
		return sl, errors.New("sleep log from fitbit: dates do not match with current day")
	}
	raw, err := json.Marshal(&s)
	if err != nil {
		return
	}

	title := "Nap"
	if s.IsMainSleep {
		title = "Main sleep"
	}
	sl = core.SleepLog{
		Date:        date,
		IsMainSleep: s.IsMainSleep,
		TimeAsleep:  fitbit.ToDuration(s.MinutesAsleep),
		TimeInBed:   fitbit.ToDuration(s.TimeInBed),
		Origin:      core.LogOriginFitbit,
		Raw:         raw,
		TimelineLog: core.TimelineLog{
			Title:     title,
			StartTime: fitbit.ParseTime(s.StartTime),
			EndTime:   fitbit.ParseTime(s.EndTime),
			Type:      core.TimelineTypeSleep,
			Note:      "Origin: " + core.LogOriginFitbit,
		},
	}
	return sl, nil
}
