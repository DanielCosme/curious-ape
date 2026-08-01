package integration

import (
	"encoding/json"
	"errors"
	"log/slog"

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
