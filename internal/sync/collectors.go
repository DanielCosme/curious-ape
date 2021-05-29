package sync

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
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

func (co *Collectors) GetTodayLog() {
	today := time.Now()
	strTime := today.Format("2006-01-02")
	co.GetLog(strTime)
}

func (co *Collectors) GetLog(date string) {
	log.Println("Getting Record")
	jsonRecord, err := co.Sleep.DayLog(date)
	if err != nil {
		log.Println("ERR on GetLog", err.Error())
	}

	err = co.saveLog(date, jsonRecord)
	if err != nil {
		log.Println("ERR on GetLog", err.Error())
	}
}

func (co *Collectors) FromDayZero(limit time.Time) error {
	zero, _ := time.Parse("2006-01-02", fitbit.ZeroDay) // From here
	err := co.RangeLogs(zero, limit)
	if err != nil {
		return err
	}
	return nil
}

func (co *Collectors) RangeLogs(zero, limit time.Time) error {
	maxYear, maxMonth, _ := limit.Date()

	first := zero
	for {
		last := first.AddDate(0, 1, -1)
		firstStr := first.Format("2006-01-02")
		lastStr := last.Format("2006-01-02")

		// Get logs for current month
		if err := co.rangeLog(firstStr, lastStr); err != nil {
			log.Println("ERR", err.Error())
			return err
		}

		currYear, currMonth, _ := first.Date()
		if currYear == maxYear && currMonth == maxMonth {
			break
		}

		first = first.AddDate(0, 1, 0)
	}

	return nil
}

// Request the record for the given months range and save them in the database on day at
// a time.
func (co *Collectors) rangeLog(first, last string) error {
	logs, err := co.Sleep.LogsRange(first, last)
	if err != nil {
		return err
	}

	// Iterate over collection and insert one record at a time.
	log.Println(len(logs), "Records To save")
	count := 0
	for k, v := range logs {
		err = co.saveLog(k, v)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE") {
				continue
			} else {
				return err
			}
		}
		count++
	}
	log.Printf("%v Values saved", count)

	return nil
}

// parse json response and save it into database
func (co *Collectors) saveLog(date string, jsonResponse []byte) error {
	var sleepRecord data.SleepRecord
	err := json.Unmarshal(jsonResponse, &sleepRecord)
	sleepRecord.RawJson = jsonResponse
	sleepRecord.Provider = co.Sleep.Auth.Provider
	if err != nil {
		log.Println("ERR", err.Error())
		return err
	}

	err = co.Models.SleepRecords.Insert(sleepRecord)
	if err != nil {
		log.Println("ERR", err.Error())
		return err
	}

	log.Println("All Good")
	return nil
}
