package application

import (
	"encoding/json"
	"errors"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/database"
	"github.com/danielcosme/curious-ape/pkg/database/gen/models"
	"github.com/danielcosme/curious-ape/pkg/integrations/fitbit"
	"time"
)

func (a *App) sleepSync(d core.Date) error {
	sls, err := a.sleepLogsGetFromFitbit(d)
	if err != nil {
		return err
	}
	for _, setter := range sls {
		sl, err := a.db.Sleep.Upsert(setter)
		if err != nil {
			return err
		}
		dur := fitbit.ToDuration(int(sl.MinutesAsleep.GetOrZero()))
		a.Log.Info("Sleep log added", "date", sl.Date, "duration", dur.String())
		if sl.IsMainSleep.GetOrZero() {
			habitState := core.HabitStateNotDone
			wakeUpTime := time.Date(sl.EndTime.Year(), sl.EndTime.Month(), sl.EndTime.Day(), 6, 0, 0, 0, sl.EndTime.Location())
			if sl.EndTime.Before(wakeUpTime) {
				habitState = core.HabitStateDone
			}
			_, err := a.HabitUpsert(d, core.HabitKindWakeUp, habitState)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *App) sleepLogsGetFromFitbit(dates ...core.Date) (res []*models.SleepLogSetter, err error) {
	fitbitClient, err := a.fitbitClient()
	if err != nil {
		return
	}
	for _, date := range dates {
		day, err := a.db.Days.GetOrCreate(database.DayParams{Date: date})
		if err != nil {
			return res, err
		}
		sleepLogs, err := fitbitClient.Sleep.GetByDate(day.Date)
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

func sleepLogFromFitbit(day *models.Day, s fitbit.Sleep) (*models.SleepLogSetter, error) {
	date := core.NewDate(day.Date)
	if !date.IsEqual(fitbit.ParseDate(s.DateOfSleep)) {
		return nil, errors.New("sleep log from fitbit: dates do not match with current day")
	}
	raw, err := json.Marshal(&s)
	if err != nil {
		return nil, err
	}
	sleepLog := &models.SleepLogSetter{
		DayID:          omit.From(day.ID),
		Date:           omit.From(day.Date),
		StartTime:      omit.From(fitbit.ParseTime(s.StartTime)),
		EndTime:        omit.From(fitbit.ParseTime(s.EndTime)),
		IsMainSleep:    omitnull.From(s.IsMainSleep),
		TotalTimeInBed: omitnull.From(int32(fitbit.ToDuration(s.TimeInBed).Minutes())),
		MinutesAsleep:  omitnull.From(int32(fitbit.ToDuration(s.MinutesAsleep).Minutes())),
		Origin:         omit.From(core.OriginLogFitbit),
		Raw:            omitnull.From(string(raw)),
	}
	return sleepLog, nil
}
