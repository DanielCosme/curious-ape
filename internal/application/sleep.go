package application

import (
	"encoding/json"
	errors2 "errors"
	"fmt"
	database2 "github.com/danielcosme/curious-ape/internal/database"
	entity2 "github.com/danielcosme/curious-ape/internal/entity"
	"time"

	"github.com/danielcosme/curious-ape/internal/integrations/fitbit"
	"github.com/danielcosme/go-sdk/errors"
	"github.com/danielcosme/go-sdk/log"
)

func (a *App) SleepDeleteByID(id int) error {
	return a.db.SleepLogs.Delete(id)
}

func (a *App) SleepUpdate(sleepLog, data *entity2.SleepLog) (*entity2.SleepLog, error) {
	if err := checkManualOrigin(sleepLog.Origin); err != nil {
		return nil, err
	}

	// Habit Upsert From Fitness Log
	if err := a.HabitUpsertFromSleepLog(*data); err != nil {
		a.Log.Error(err)
	}

	data.ID = sleepLog.ID
	return a.db.SleepLogs.Update(data, database2.SleepLogsJoinDay(a.db))
}

func (a *App) SleepCreateFromApi(sleepLog *entity2.SleepLog) (*entity2.SleepLog, error) {
	if sleepLog.Raw == "" {
		start := sleepLog.StartTime.AddDate(0, 0, -1)
		sleepLog.MinutesInBed = sleepLog.EndTime.Sub(start)
		sleepLog.MinutesAsleep = time.Duration(float64(sleepLog.MinutesInBed) * 0.87)
		// Calculate an average of 13% of bedtime awake
		sleepLog.MinutesAwake = time.Duration(float64(sleepLog.MinutesInBed) * 0.13)
		a.Log.DebugP("manual sleep log", log.Prop{
			"Time in bed": sleepLog.MinutesInBed.String(),
			"Asleep":      sleepLog.MinutesAwake.String(),
			"Awake":       sleepLog.MinutesAwake.String(),
		})
	}

	if err := a.createSleepLogs([]*entity2.SleepLog{sleepLog}); err != nil {
		return nil, err
	}
	return a.db.SleepLogs.Get(entity2.SleepLogFilter{ID: []int{sleepLog.ID}}, database2.SleepLogsJoinDay(a.db))
}

func (a *App) SleepLogGet(f entity2.SleepLogFilter) (*entity2.SleepLog, error) {
	return a.db.SleepLogs.Get(f, database2.SleepLogsJoinDay(a.db))
}

func (a *App) SleepLogsFind(f entity2.SleepLogFilter) ([]*entity2.SleepLog, error) {
	return a.db.SleepLogs.Find(f, database2.SleepLogsJoinDay(a.db))
}

func (a *App) SyncSleep() error {
	days, err := a.db.Days.Find(entity2.DayFilter{}, database2.DaysJoinSleepLogs(a.db))
	if err != nil {
		return err
	}
	client, err := a.Oauth2GetClient(entity2.ProviderFitbit)
	if err != nil {
		return err
	}
	fitbitApi := a.sync.FitbitClient(client)

	for _, d := range days {
		// TODO this will not work properly for more than one provider at a time.
		if len(d.SleepLogs) == 0 {
			// Try to sync if we don't have the sleep log
			payload, err := fitbitApi.Sleep.GetByDate(d.Date)
			if err != nil {
				return err
			}
			sleepLogs, err := toSleepLogFromFitbit([]*entity2.Day{d}, payload.Sleep)
			if err != nil {
				return err
			}
			// Save log
			if err := a.createSleepLogs(sleepLogs); err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *App) SyncSleepByDateRange(start, end time.Time) error {
	// Get client, refreshes token if necessary
	client, err := a.Oauth2GetClient(entity2.ProviderFitbit)
	if err != nil {
		return err
	}

	// Get sleep records from fitbit
	fitbitPayload, err := a.sync.FitbitClient(client).Sleep.GetByDateRange(start, end)
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

	return a.createSleepLogs(sleepLogs)
}

func (a *App) SyncFitbitSleepLog(date time.Time) error {
	day, err := database2.DayGetOrCreate(a.db, date)
	if err != nil {
		return err
	}

	// Get client, refreshes token if necessary
	client, err := a.Oauth2GetClient(entity2.ProviderFitbit)
	if err != nil {
		return err
	}

	// Get sleep records from fitbit
	fitbitPayload, err := a.sync.FitbitClient(client).Sleep.GetByDate(date)
	if err != nil {
		return err
	}
	// Map fitbit records to application sleep log struct
	sleepLogs, err := toSleepLogFromFitbit([]*entity2.Day{day}, fitbitPayload.Sleep)
	if err != nil {
		return err
	}

	// Persist
	return a.createSleepLogs(sleepLogs)
}

func toSleepLogFromFitbit(days []*entity2.Day, sleepRecords []fitbit.Sleep) ([]*entity2.SleepLog, error) {
	sleepLogs := make([]*entity2.SleepLog, 0, len(sleepRecords))
	mapDays := database2.DayToMapByISODate(days)

	for _, s := range sleepRecords {
		var day *entity2.Day
		if d, ok := mapDays[s.DateOfSleep]; ok {
			day = d
		} else {
			return nil, errors.New(fmt.Sprintf("not day match for fitbit log on %s", s.DateOfSleep))
		}

		sleepLog := &entity2.SleepLog{
			Day:           day,
			DayID:         day.ID,
			Date:          fitbit.ParseDate(s.DateOfSleep),
			StartTime:     fitbit.ParseTime(s.StartTime),
			EndTime:       fitbit.ParseTime(s.EndTime),
			IsMainSleep:   s.IsMainSleep,
			IsAutomated:   true,
			Origin:        entity2.Fitbit,
			MinutesInBed:  fitbit.ToDuration(s.TimeInBed),
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

func (a *App) createSleepLogs(logs []*entity2.SleepLog) error {
	for _, l := range logs {
		op := "Created"
		testerLog, err := a.db.SleepLogs.Get(entity2.SleepLogFilter{DayID: []int{l.Day.ID}})
		if err != nil && !errors.Is(err, database2.ErrNotFound) {
			return err
		}
		if testerLog != nil && testerLog.StartTime.Equal(l.StartTime) && testerLog.EndTime.Equal(l.EndTime) {
			// update
			op = "Updated"
			l.ID = testerLog.ID
			if _, err := a.db.SleepLogs.Update(l); err != nil {
				return fmt.Errorf("sleep log could not be updated: %w", err)
			}
		} else if err := a.db.SleepLogs.Create(l); err != nil {
			if errors2.Is(err, database2.ErrUniqueCheckFailed) {
				a.Log.Error(fmt.Errorf("dayID and main sleep unique check failed for %s", l.Date.Format(entity2.HumanDateWithTime)))
			}
			return err
		}

		// Habit Creation From Fitness Log
		if err := a.HabitUpsertFromSleepLog(*l); err != nil {
			a.Log.Error(err)
		}

		a.Log.InfoP(fmt.Sprintf("%s sleep log", op), log.Prop{
			"provider": l.Origin.Str(),
			"date":     l.Date.Format(entity2.HumanDateWithTime),
		})
	}

	return nil
}
