package application

import (
	fitbit2 "github.com/danielcosme/curious-ape/fitbit/fitbit"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/datasource"
	"time"
)

func (a *App) SleepDebug() (*entity.SleepLog, error) {
	// Get client, refreshes token if necessary
	client, err := a.Oauth2GetClient(entity.ProviderFitbit)
	if err != nil {
		return nil, err
	}

	api := fitbit2.NewAPI(client)
	fitbitSl, err := api.Sleep.GetLogByDate(time.Now())
	log := fitbitSl.Sleep[0]

	day, err := a.DayGetByDate(parseFitbitDate(log.DateOfSleep))
	if err != nil {
		return nil, err
	}

	mySl := &entity.SleepLog{
		DayID:          day.ID,
		StartTime:      parseFitbitTime(log.StartTime),
		EndTime:        parseFitbitTime(log.EndTime),
		IsMainSleep:    log.IsMainSleep,
		IsAutomated:    true,
		Origin:         "fitbit",
		TotalTimeInBed: log.TimeInBed,
		MinutesAsleep:  log.TimeInBed,
		MinutesRem:     fitbitSl.Summary.Stages.Rem,
		MinutesDeep:    fitbitSl.Summary.Stages.Deep,
		MinutesLight:   fitbitSl.Summary.Stages.Light,
		MinutesAwake:   log.MinutesAwake,
	}

	if err := a.db.SleepLogs.Create(mySl); err != nil {
		return nil, err
	}

	return a.db.SleepLogs.Get(entity.SleepLogFilter{ID: []int{mySl.ID}}, datasource.SleepLogsPipeline(a.db)...)
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
