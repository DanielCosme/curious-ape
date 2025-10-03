package application

import (
	"context"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/oak"
)

func (a *App) sleepSync(ctx context.Context, d core.Date) error {
	logger := oak.FromContext(ctx)
	logger.Notice("Sleep Sync not implemented")

	/*
		sls, err := a.sleepLogsGetFromFitbit(d)
		if err != nil {
			return err
		}
		for _, setter := range sls {
			sl, err := a.db.Sleep.Upsert(setter)
			if err != nil {
				return err
			}
			dur := fitbit.ToDuration(int(sl.MinutesAsleep))
			logger.Info("Sleep log added", "date", sl.Date, "duration", dur.String())
			if sl.IsMainSleep {
				habitState := core.HabitStateNotDone
				wakeUpTime := time.Date(sl.EndTime.Year(), sl.EndTime.Month(), sl.EndTime.Day(), 6, 0, 0, 0, sl.EndTime.Location())
				if sl.EndTime.Before(wakeUpTime) {
					habitState = core.HabitStateDone
				}
				_, err := a.HabitUpsert(ctx, core.UpsertHabitParams{
					Date:      d,
					Type:      core.HabitTypeWakeUp,
					State:     habitState,
					Automated: true})
				if err != nil {
					return err
				}
			}
		}
	*/
	return nil
}

/*
func (a *App) sleepLogsGetFromFitbit(dates ...core.Date) (res []*models.SleepLogSetter, err error) {
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

func sleepLogFromFitbit(day core.Day, s fitbit.Sleep) (*models.SleepLogSetter, error) {
	if !day.Date.IsEqual(fitbit.ParseDate(s.DateOfSleep)) {
		return nil, errors.New("sleep log from fitbit: dates do not match with current day")
	}
	raw, err := json.Marshal(&s)
	if err != nil {
		return nil, err
	}
	sleepLog := &models.SleepLogSetter{
		DayID:         omit.From(int64(day.ID)),
		Date:          omit.From(day.Date.Time()),
		StartTime:     omit.From(fitbit.ParseTime(s.StartTime)),
		EndTime:       omit.From(fitbit.ParseTime(s.EndTime)),
		IsMainSleep:   omit.From(s.IsMainSleep),
		MinutesAsleep: omit.From(int64(fitbit.ToDuration(s.MinutesAsleep).Minutes())),
		Origin:        omit.From(core.OriginLogFitbit),
		Raw:           omit.From(string(raw)),
	}
	return sleepLog, nil
}
*/
