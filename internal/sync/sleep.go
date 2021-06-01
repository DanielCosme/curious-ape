package sync

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/danielcosme/curious-ape/internal/data"
	"github.com/danielcosme/curious-ape/internal/sync/fitbit"
)

type SleepCollector struct {
	Models *data.Models
	*fitbit.SleepProvider
}

func (co *SleepCollector) GetCurrentDayRecord() error {
	strTime := time.Now().Format("2006-01-02")
	err := co.GetRecord(strTime)
	if err != nil {
		return err
	}
	return nil
}

func (co *SleepCollector) GetRecord(date string) error {
	log.Println("Getting Record")
	jsonRecord, err := co.DayLog(date)
	if err != nil {
		log.Println("ERR on GetLog", err.Error())
		return err
	}

	err = co.saveLog(jsonRecord)
	if err != nil {
		log.Println("ERR on GetLog", err.Error())
		return err
	}

	return nil
}

func (co *SleepCollector) GetRecordsFromDayZero(limit time.Time) error {
	zero, _ := time.Parse("2006-01-02", fitbit.ZeroDay) // From here
	err := co.GetRecords(zero, limit)
	if err != nil {
		return err
	}
	return nil
}

func (co *SleepCollector) GetRecords(zero, limit time.Time) error {
	maxYear, maxMonth, _ := limit.Date()

	first := zero
	for {
		last := first.AddDate(0, 1, -1)
		firstStr := first.Format("2006-01-02")
		lastStr := last.Format("2006-01-02")

		// Get logs for current month
		if err := co.getMonthRecords(firstStr, lastStr); err != nil {
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
func (co *SleepCollector) getMonthRecords(first, last string) error {
	logs, err := co.LogsRange(first, last)
	if err != nil {
		return err
	}

	// Iterate over collection and insert one record at a time.
	log.Println(len(logs), "Records To save")
	count := 0
	for _, v := range logs {
		err = co.saveLog(v)
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
func (co *SleepCollector) saveLog(jsonResponse []byte) error {
	var sleepRecord data.SleepRecord
	err := json.Unmarshal(jsonResponse, &sleepRecord)
	sleepRecord.RawJson = jsonResponse
	sleepRecord.Provider = co.Auth.Provider
	if err != nil {
		log.Println("ERR", err.Error())
		return err
	}

	err = co.Models.SleepRecords.Insert(sleepRecord)
	if err != nil {
		log.Println("ERR", err.Error())
		return err
	}

	err = co.saveSleepHabit(&sleepRecord)
	if err != nil {
		log.Println("ERR", err.Error())
		return err
	}

	log.Println("All Good")
	return nil
}

func (co *SleepCollector) saveSleepHabit(sleepRecord *data.SleepRecord) error {
	var habit *data.Habit = &data.Habit{
		Date:   sleepRecord.Date,
		Origin: sleepRecord.Provider,
		Type:   "sleep",
		State:  "yes",
	}

	d := strings.Split(sleepRecord.EndTime, "T")
	hour := d[1][:2]
	hourInt, err := strconv.Atoi(hour)
	if err != nil {
		return err
	}

	if hourInt > 6 {
		habit.State = "no"
	}

	err = co.Models.Habits.UpdateOrCreate(habit)
	if err != nil {
		return err
	}

	return nil
}

func (co *SleepCollector) BuildHabitsFromSleepRecords() (err error) {
	all, err := co.Models.SleepRecords.GetAll()
	if err != nil {
		return err
	}

	l := len(all)
	for i, v := range all {
		err := co.saveSleepHabit(v)
		if err != nil {
			return err
		}

		if i+1 == l {
			log.Println("succesfully added", l, "habits")
		}
	}

	return nil
}
