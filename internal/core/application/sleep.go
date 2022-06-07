package application

import (
	"github.com/danielcosme/curious-ape/fitbit/fitbit"
	"github.com/danielcosme/curious-ape/internal/core/database"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"time"
)

func (a *App) GetSleepLogsForDay(d *entity.Day) ([]*entity.SleepLog, error) {
	// Get client, refreshes token if necessary
	client, err := a.Oauth2GetClient(entity.ProviderFitbit)
	if err != nil {
		return nil, err
	}

	api := fitbit.NewAPI(client)
	fitbitSl, err := api.Sleep.GetLogByDate(d.Date)
	if err != nil {
		return nil, err
	}

	logs := FromFitbitSleepToLog(d, fitbitSl)
	for _, l := range logs {
		if err := a.db.SleepLogs.Create(l); err != nil {
			return nil, err
		}
	}

	return a.db.SleepLogs.Find(entity.SleepLogFilter{ID: database.SleepToIDs(logs)}, database.SleepLogsPipeline(a.db)...)
}

func FromFitbitSleepToLog(d *entity.Day, sleepEnvelope *fitbit.SleepEnvelope) []*entity.SleepLog {
	sleepLogs := []*entity.SleepLog{}

	for _, fsl := range sleepEnvelope.Sleep {
		mySl := &entity.SleepLog{
			DayID:         d.ID,
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
		sleepLogs = append(sleepLogs, mySl)
	}

	return sleepLogs
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
