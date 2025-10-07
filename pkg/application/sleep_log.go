package application

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/integrations/fitbit"
	"github.com/danielcosme/curious-ape/pkg/oak"
)

func (a *App) sleepSync(ctx context.Context, d core.Date) error {
	logger := oak.FromContext(ctx)

	sls, err := a.sleepLogsGetFromFitbit(d)
	if err != nil {
		return err
	}
	for _, slParams := range sls {
		sl, err := a.db.Sleep.Upsert(slParams)
		if err != nil {
			return err
		}
		logger.Info("Sleep log added", "date", sl.Date, "duration", sl.TimeAsleep)
		if sl.IsMainSleep {
			habitState := core.HabitStateNotDone
			wakeUpTime := time.Date(sl.EndTime.Year(), sl.EndTime.Month(), sl.EndTime.Day(), 6, 0, 0, 0, sl.EndTime.Location())
			if sl.EndTime.Before(wakeUpTime) {
				habitState = core.HabitStateDone
			}
			_, err := a.HabitUpsert(ctx, core.Habit{
				Date:      d,
				Type:      core.HabitTypeWakeUp,
				State:     habitState,
				Note:      sl.EndTime.Format(core.Time),
				Automated: true})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *App) sleepLogsGetFromFitbit(dates ...core.Date) (res []core.SleepLog, err error) {
	fitbitClient, err := a.fitbitClient()
	if err != nil {
		return
	}
	for _, date := range dates {
		day, err := a.dayGetOrCreate(date)
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

func sleepLogFromFitbit(day core.Day, s fitbit.Sleep) (sl core.SleepLog, err error) {
	if !day.Date.IsEqual(fitbit.ParseDate(s.DateOfSleep)) {
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
		Date:        day.Date,
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
