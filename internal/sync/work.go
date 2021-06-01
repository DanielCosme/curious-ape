package sync

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/danielcosme/curious-ape/internal/data"
	"github.com/danielcosme/curious-ape/internal/sync/toggl"
)

type WorkCollector struct {
	Models *data.Models
	*toggl.WorkProvider
}

func (co *WorkCollector) GetRecords(start, end string) error {
	log.Println("START", start)
	result, err := co.LogsRange(start, end)
	if err != nil {
		return err
	}

	c := 0
	l := len(result)
	log.Println(l, "Work Records to save")
	for k, v := range result {

		err := co.decodeAndSave(k, v)
		if err != nil {
			return err
		}
		c++
	}
	log.Println(c, "Work Records saved")

	return nil
}

func (co *WorkCollector) GetRecordsFromDayZero(limit time.Time) error {
	zero := "2021-01-01"
	log.Println(zero, limit)
	err := co.GetRecords(zero, limit.Format("2006-01-02"))
	if err != nil {
		return err
	}
	return nil
}

func (co *WorkCollector) GetCurrentDayRecord() error {
	strTime := time.Now().Format("2006-01-02")
	err := co.GetRecord(strTime)
	if err != nil {
		return err
	}
	return nil
}

func (co *WorkCollector) GetRecord(date string) error {
	log.Println("Getting Record")
	jsonRecord, err := co.DayLog(date)
	if err != nil {
		log.Println("ERR on GetLog", err.Error())
		return err
	}

	err = co.decodeAndSave(date, jsonRecord)
	if err != nil {
		return err
	}

	return nil
}

func (co *WorkCollector) decodeAndSave(date string, jsonRecord []byte) error {
	r, err := co.decode(date, jsonRecord)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			log.Println(err.Error(), "for", date)
			// log habit as no, not done or log it into the: to manually revise list.
			return nil
		}
		return err
	}

	err = co.saveLog(r)
	if err != nil {
		log.Println("ERR on GetLog", err.Error())
		return err
	}

	return nil
}

func (co *WorkCollector) decode(date string, jsonResponse []byte) (*data.WorkRecord, error) {
	var jsonMap map[string]interface{}
	workRecord := &data.WorkRecord{}

	if err := json.Unmarshal(jsonResponse, &jsonMap); err != nil {
		return nil, err
	}
	total, ok := jsonMap["total_grand"].(float64)
	if !ok {
		return nil, ErrNoRecord
	}

	workRecord.Date = date
	workRecord.Total = int(total)
	workRecord.RawJson = string(jsonResponse)
	workRecord.Provider = co.Auth.Provider

	return workRecord, nil
}

func (co *WorkCollector) saveLog(wr *data.WorkRecord) error {
	err := co.Models.WorkRecords.Insert(wr)
	if err != nil {
		return err
	}

	err = co.saveWorkHabit(wr)
	if err != nil {
		return err
	}

	return nil
}

func (co *WorkCollector) saveWorkHabit(wr *data.WorkRecord) error {
	log.Println("Saving Work habit for", wr.Date)
	var habit *data.Habit = &data.Habit{
		Date:   wr.Date,
		Origin: wr.Provider,
		Type:   "work",
		State:  "yes",
	}

	targetInMiliseconds := 18000000 // 5 hours
	if wr.Total < targetInMiliseconds {
		habit.State = "no"
	}

	err := co.Models.Habits.UpdateOrCreate(habit)
	if err != nil {
		return err
	}

	return nil
}

func (co *WorkCollector) BuildHabitsFromWorkRecords() (err error) {
	all, err := co.Models.WorkRecords.GetAll()
	if err != nil {
		return err
	}

	l := len(all)
	for i, v := range all {
		err := co.saveWorkHabit(v)
		if err != nil {
			return err
		}

		if i+1 == l {
			log.Println("succesfully added", l, "habits")
		}
	}

	return nil
}

func milisecondsToHours(mil int) int {
	return 0
}
