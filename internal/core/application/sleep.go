package application

import (
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
	// Get client, refreshes token if necessary
	// client, err := a.Oauth2GetClient(entity.ProviderFitbit)
	// if err != nil {
	// 	return nil, err
	// }

	// fitbitSl, err := a.Sync.FitbitClient(client).Sleep.GetLogByDate(d.Date)
	// if err != nil {
	// 	return nil, err
	// }

	// logs := fromFitbitSleepToLog(d, fitbitSl)
	// for _, l := range logs {
	// 	if err := a.db.SleepLogs.Create(l); err != nil {
	// 		return nil, err
	// 	}
	// }

	return a.getSleepLogs(entity.SleepLogFilter{DayID: []int{d.ID}})
}

func (a *App) GetAllSleepLogs() ([]*entity.SleepLog, error) {
	return a.getSleepLogs(entity.SleepLogFilter{})
}

func (a *App) getSleepLogs(f entity.SleepLogFilter) ([]*entity.SleepLog, error) {
	return a.db.SleepLogs.Find(f, database.SleepLogsJoinDay(a.db))
}

func (a *App) SyncSleepLogs(start, end time.Time) error {
	// Get client, refreshes token if necessary
	client, err := a.Oauth2GetClient(entity.ProviderFitbit)
	if err != nil {
		return err
	}

	fitbitSl, err := a.Sync.FitbitClient(client).Sleep.GetLogByDateRange(start, end)
	if err != nil {
		return err
	}

	days, err := a.daysGetByDateRange(start, end)
	if err != nil {
		return err
	}

	return a.saveSleepLogsFromFitbit(days, fitbitSl)
}

func (a *App) SyncSleepLog(date time.Time) error {
	day, err := a.DayGetByDate(date)
	if err != nil {
		return err
	}

	// Get client, refreshes token if necessary
	client, err := a.Oauth2GetClient(entity.ProviderFitbit)
	if err != nil {
		return err
	}

	fitbitSl, err := a.Sync.FitbitClient(client).Sleep.GetLogByDate(date)
	if err != nil {
		return err
	}

	return a.saveSleepLogsFromFitbit([]*entity.Day{day}, fitbitSl)
}

func (a *App) saveSleepLogsFromFitbit(days []*entity.Day, sleepEnvelope *fitbit.SleepEnvelope) error {
	mapDays := database.DayToMapByISODate(days)

	for _, fsl := range sleepEnvelope.Sleep {
		dayID := 0
		if d, ok := mapDays[fsl.DateOfSleep]; ok {
			dayID = d.ID
		} else {
			return errors.New(fmt.Sprintf("not day match for fitbit log on %s", fsl.DateOfSleep))
		}

		mySl := &entity.SleepLog{
			DayID:         dayID,
			StartTime:     parseFitbitTime(fsl.StartTime),
			EndTime:       parseFitbitTime(fsl.EndTime),
			IsMainSleep:   fsl.IsMainSleep,
			IsAutomated:   true,
			Origin:        "fitbit",
			TimeInBed:     fromIntMinutesToDuration(fsl.TimeInBed),
			MinutesAsleep: fromIntMinutesToDuration(fsl.MinutesAsleep),
			MinutesAwake:  fromIntMinutesToDuration(fsl.MinutesAwake),
			MinutesRem:    fromIntMinutesToDuration(fsl.Levels.Summary.Rem.Minutes),
			MinutesDeep:   fromIntMinutesToDuration(fsl.Levels.Summary.Deep.Minutes),
			MinutesLight:  fromIntMinutesToDuration(fsl.Levels.Summary.Light.Minutes),
		}

		if err := a.db.SleepLogs.Create(mySl); err != nil {
			if errors2.Is(err, database.ErrUniqueCheckFailed) {
				a.Log.Error(err)
			} else {
				return err
			}
		} else {
			prop := log.Prop{
				"provider": "fitbit",
				"date": fsl.DateOfSleep,
			}
			a.Log.InfoP("Created sleep log", prop)
		}
	}

	return nil
}

func fromIntMinutesToDuration(i int) time.Duration {
	return time.Duration(i) * time.Minute
}

func parseFitbitDate(s string) time.Time {
	// yyyy-mm-dd
	t, _ := time.Parse("2006-01-02", s)
	return t
}

func parseFitbitTime(s string) time.Time {
	// 2022-06-02T05:18:30.000
	t, _ := time.Parse("2006-01-02T15:04:05.999", s)
	return t
}
