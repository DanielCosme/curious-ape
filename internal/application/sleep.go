package application

import (
	"encoding/json"
	"errors"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database"
	"github.com/danielcosme/curious-ape/internal/integrations/fitbit"
)

func (a *App) sleepSync(d core.Date) error {
	sls, err := a.sleepLogsGetFromFitbit(d)
	if err != nil {
		return err
	}
	for _, sl := range sls {
		if sl.IsMainSleep {
			habitLogParams := sl.ToHabitLogWakeUp()
			_, err := a.HabitUpsert(habitLogParams)
			if err != nil {
				return err
			}
		}
		sl, err = a.db.Sleep.Upsert(sl)
		if err != nil {
			return err
		}
		a.Log.Info("Sleep log added", "date", sl.Date, "duration", sl.MinutesAsleep.String())
	}
	return nil
}

func (a *App) sleepLogsGetFromFitbit(dates ...core.Date) (res []core.SleepLog, err error) {
	fitbitClient, err := a.fitbitClient()
	if err != nil {
		return
	}
	for _, date := range dates {
		day, err := a.db.Days.GetOrCreate(database.DayParams{Date: date})
		if err != nil {
			return res, err
		}
		sleepLogs, err := fitbitClient.Sleep.GetByDate(day.Date.Time())
		if err != nil {
			return res, err
		}
		for _, fsl := range sleepLogs.Sleep {
			sl, err := sleepLogFromFitbit(day, fsl)
			if err != nil {
				return res, err
			}
			res = append(res, sl)
		}
	}
	return
}

func sleepLogFromFitbit(day core.Day, s fitbit.Sleep) (core.SleepLog, error) {
	sleepLog := core.NewSleepLog(day)
	if !day.Date.IsEqual(fitbit.ParseDate(s.DateOfSleep)) {
		return core.SleepLog{}, errors.New("sleep log from fitbit: dates do not match with current day")
	}
	raw, err := json.Marshal(&s)
	if err != nil {
		return core.SleepLog{}, err
	}
	sleepLog.Raw = raw
	sleepLog.IsAutomated = true
	sleepLog.Origin = core.IntegrationFitbit
	sleepLog.IsMainSleep = s.IsMainSleep
	sleepLog.StartTime = fitbit.ParseTime(s.StartTime)
	sleepLog.EndTime = fitbit.ParseTime(s.EndTime)
	sleepLog.MinutesInBed = fitbit.ToDuration(s.TimeInBed)
	sleepLog.MinutesAsleep = fitbit.ToDuration(s.MinutesAsleep)
	sleepLog.MinutesAwake = fitbit.ToDuration(s.MinutesAwake)
	sleepLog.MinutesRem = fitbit.ToDuration(s.Levels.Summary.Rem.Minutes)
	sleepLog.MinutesDeep = fitbit.ToDuration(s.Levels.Summary.Deep.Minutes)
	sleepLog.MinutesLight = fitbit.ToDuration(s.Levels.Summary.Light.Minutes)
	return sleepLog, nil
}
