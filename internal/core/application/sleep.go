package application

import (
	"encoding/json"
	errors2 "errors"
	"fmt"
	"github.com/danielcosme/curious-ape/internal/core/database"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/integrations/fitbit"
	"github.com/danielcosme/curious-ape/sdk/errors"
	"github.com/danielcosme/curious-ape/sdk/log"
	"time"
)

func (a *App) GetSleepLogsForDay(d *entity.Day) ([]*entity.SleepLog, error) {
	return a.getSleepLogs(entity.SleepLogFilter{DayID: []int{d.ID}})
}

func (a *App) GetAllSleepLogs() ([]*entity.SleepLog, error) {
	return a.getSleepLogs(entity.SleepLogFilter{})
}

func (a *App) getSleepLogs(f entity.SleepLogFilter) ([]*entity.SleepLog, error) {
	return a.db.SleepLogs.Find(f, database.SleepLogsJoinDay(a.db))
}

func (a *App) FitbitSyncSleepLogs() error {
	days, err := a.db.Days.Find(entity.DayFilter{}, database.DaysJoinSleepLogs(a.db))
	if err != nil {
		return err
	}
	client, err := a.Oauth2GetClient(entity.ProviderFitbit)
	if err != nil {
		return err
	}
	fitbitApi := a.Sync.FitbitClient(client)

	for _, d := range days {
		if len(d.SleepLogs) == 0 {
			// Try to sync if we don't have the sleep log
			payload, err := fitbitApi.Sleep.GetByDate(d.Date)
			if err != nil {
				return err
			}
			sleepLogs, err := toSleepLogFromFitbit([]*entity.Day{d}, payload.Sleep)
			if err != nil {
				return err
			}
			// Save log
			if err := a.SaveSleepLogs(sleepLogs); err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *App) FitbitSyncSleepLogsByDateRange(start, end time.Time) error {
	// Get client, refreshes token if necessary
	client, err := a.Oauth2GetClient(entity.ProviderFitbit)
	if err != nil {
		return err
	}

	// Get sleep records from fitbit
	fitbitPayload, err := a.Sync.FitbitClient(client).Sleep.GetByDateRange(start, end)
	if err != nil {
		return err
	}
	days, err := a.daysGetByDateRange(start, end)
	if err != nil {
		return err
	}
	// Map fitbit records to application sleep log struct
	sleepLogs, err := toSleepLogFromFitbit(days, fitbitPayload.Sleep)
	if err != nil {
		return err
	}

	return a.SaveSleepLogs(sleepLogs)
}

func (a *App) FitbitSyncSleepLog(date time.Time) error {
	day, err := a.DayGetByDate(date)
	if err != nil {
		return err
	}

	// Get client, refreshes token if necessary
	client, err := a.Oauth2GetClient(entity.ProviderFitbit)
	if err != nil {
		return err
	}

	// Get sleep records from fitbit
	fitbitPayload, err := a.Sync.FitbitClient(client).Sleep.GetByDate(date)
	if err != nil {
		return err
	}
	// Map fitbit records to application sleep log struct
	sleepLogs, err := toSleepLogFromFitbit([]*entity.Day{day}, fitbitPayload.Sleep)
	if err != nil {
		return err
	}

	// Persist
	return a.SaveSleepLogs(sleepLogs)
}

func toSleepLogFromFitbit(days []*entity.Day, sleepRecords []fitbit.Sleep) ([]*entity.SleepLog, error) {
	sleepLogs := make([]*entity.SleepLog, 0, len(sleepRecords))
	mapDays := database.DayToMapByISODate(days)

	for _, s := range sleepRecords {
		var day *entity.Day
		if d, ok := mapDays[s.DateOfSleep]; ok {
			day = d
		} else {
			return nil, errors.New(fmt.Sprintf("not day match for fitbit log on %s", s.DateOfSleep))
		}

		sleepLog := &entity.SleepLog{
			Day:           day,
			DayID:         day.ID,
			Date:          fitbit.ParseDate(s.DateOfSleep),
			StartTime:     fitbit.ParseTime(s.StartTime),
			EndTime:       fitbit.ParseTime(s.EndTime),
			IsMainSleep:   s.IsMainSleep,
			IsAutomated:   true,
			Origin:        entity.Fitbit,
			TimeInBed:     fitbit.ToDuration(s.TimeInBed),
			MinutesAsleep: fitbit.ToDuration(s.MinutesAsleep),
			MinutesAwake:  fitbit.ToDuration(s.MinutesAwake),
			MinutesRem:    fitbit.ToDuration(s.Levels.Summary.Rem.Minutes),
			MinutesDeep:   fitbit.ToDuration(s.Levels.Summary.Deep.Minutes),
			MinutesLight:  fitbit.ToDuration(s.Levels.Summary.Light.Minutes),
		}

		raw, err := json.Marshal(s)
		if err != nil {
			return nil, err
		}
		sleepLog.Raw = string(raw) // save raw json

		sleepLogs = append(sleepLogs, sleepLog)
	}

	return sleepLogs, nil
}

func (a *App) SaveSleepLogs(logs []*entity.SleepLog) error {
	for _, l := range logs {
		testerLog, err := a.db.SleepLogs.Get(entity.SleepLogFilter{DayID: []int{l.Day.ID}})
		if err != nil && !errors.Is(err, database.ErrNotFound) {
			return err
		}
		if testerLog != nil && testerLog.StartTime.Equal(l.StartTime) && testerLog.EndTime.Equal(l.EndTime) {
			a.Log.Warningf("Sleep log for %s already exist, not going to be saved.", l.Date.Format(entity.HumanDate))
			continue
		}

		// Habit Creation From Sleep Log
		if err := a.HabitCreateFromSleepLog(*l); err != nil {
			a.Log.Error(err)
		}

		// TODO implement upsert here
		if err := a.db.SleepLogs.Create(l); err != nil {
			if errors2.Is(err, database.ErrUniqueCheckFailed) {
				a.Log.Error(fmt.Errorf("dayID and main sleep unique check failed for %s", l.Date.Format(entity.HumanDate)))
			} else {
				return err
			}
		} else {
			prop := log.Prop{
				"provider": l.Origin.Str(),
				"date":     l.Date.Format(entity.HumanDate),
			}
			a.Log.InfoP("Created sleep log", prop)
		}
	}

	return nil
}
