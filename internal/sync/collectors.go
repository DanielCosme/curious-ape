package sync

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/danielcosme/curious-ape/internal/data"
	"github.com/danielcosme/curious-ape/internal/sync/fitbit"
)

var (
	ErrTokenExpired = errors.New("token expired")
	ErrUnauthorized = errors.New("server needs to authorize again")
)

type Collectors struct {
	Models *data.Models
	Sleep  *fitbit.SleepCollector
}

func NewCollectors(models *data.Models) *Collectors {
	f := &fitbit.SleepCollector{
		Auth:  fitbit.FitbitAuth,
		Token: &models.Tokens,
		Scope: "sleep",
	}

	return &Collectors{
		Models: models,
		Sleep:  f,
	}
}

func (co *Collectors) DayLog() {
	today := time.Now()
	strTime := today.Format("2006-01-02")
	co.GetDay(strTime)
}

func (co *Collectors) GetDay(date string) {
	jsonRecord, err := co.Sleep.DayLog(date)
	if err != nil {
		log.Println("ERR", err.Error())
		return
	}

	var sleepRecord data.SleepRecord
	err = json.Unmarshal(jsonRecord, &sleepRecord)
	sleepRecord.RawJson = jsonRecord
	sleepRecord.Provider = co.Sleep.Auth.Provider
	if err != nil {
		log.Println("ERR", err.Error())
		return
	}

	err = co.Models.SleepRecords.Insert(sleepRecord)
	if err != nil {
		log.Println("ERR", err.Error())
		return
	}

	log.Println("Record Added")
}
