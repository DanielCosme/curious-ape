package sync

import (
	"github.com/danielcosme/curious-ape/internal/core"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/danielcosme/curious-ape/internal/data"
	"github.com/danielcosme/curious-ape/internal/sync/google"
)

type FitnessCollector struct {
	Models *data.Models
	*google.FitnessProvider
}

func (co *FitnessCollector) GetCurrentDayRecord() error {
	strTime := time.Now().Format("2006-01-02")
	err := co.GetRecord(strTime)
	if err != nil {
		return err
	}
	return nil
}

func (co *FitnessCollector) GetRecordsFromDayZero(limit time.Time) error {
	err := co.GetRecords("2021-01-01", limit.Format("2006-01-02"))
	if err != nil {
		return err
	}
	return nil
}

func (co *FitnessCollector) GetRecords(start, finish string) error {
	records, err := co.LogsRange(start, finish)
	if err != nil {
		return err
	}

	for _, record := range records {
		err := co.saveLog(record)
		if err != nil {
			return err
		}
	}

	err = co.BuildHabitsFromFitnessRecords()
	if err != nil {
		return err
	}

	return nil
}

func (co *FitnessCollector) GetRecord(date string) error {
	res, err := co.DayLog(date)
	if err != nil {
		if err == google.ErrNoFitRecord {
			co.saveFitnessHabit(date, "google/strong", false)
		}
		return err
	}

	err = co.saveLog(res)
	if err != nil {
		return err
	}

	co.saveFitnessHabit(date, "google/strong", true)
	return nil
}

func (co *FitnessCollector) saveFitnessHabit(date, prov string, state bool) error {
	habit := &core.Habit{
		Date:   date,
		Origin: prov,
		Type:   "fitness",
		State:  "yes",
	}

	if !state {
		habit.State = "no"
	}

	err := co.Models.Habits.UpdateOrCreate(habit)
	if err != nil {
		return err
	}
	return nil
}

func (co *FitnessCollector) saveLog(record map[string]string) error {
	start, err := strconv.Atoi(record["startTimeMillis"])
	end, err := strconv.Atoi(record["endTimeMillis"])
	if err != nil {
		return err
	}
	date := time.Unix(int64(start/1000), 0).Format("2006-01-02")
	r := &core.FitnessRecord{
		Date:               date,
		StartInMilliseconds: start,
		EndInMilliseconds:   end,
		Provider:           record["packageName"],
	}

	err = co.Models.FitnessRecords.Insert(r)
	if err != nil {
		log.Println(err)
		if strings.Contains(err.Error(), "unique constraint") {
			return nil
		}
		return err
	}

	return nil
}

func (co *FitnessCollector) BuildHabitsFromFitnessRecords() (err error) {
	// from 2021-01-01 to today, create list of dates.
	all, err := co.Models.FitnessRecords.GetAll()
	if err != nil {
		return err
	}

	recordsMap := mapRecords(all)
	dates, err := createDates()
	if err != nil {
		return err
	}

	habit := false
	for _, date := range dates {
		if _, ok := recordsMap[date]; ok {
			habit = true
		} else {
			habit = false
		}
		co.saveFitnessHabit(date, "google/strong", habit)
	}

	return nil
}

func mapRecords(col []*core.FitnessRecord) map[string]string {
	mapRecords := make(map[string]string)
	for _, v := range col {
		date := strings.Split(v.Date, "T")[0]
		mapRecords[date] = date
	}
	return mapRecords
}

func createDates() ([]string, error) {
	dates := []string{}
	maxYear, maxMonth, maxDay := time.Now().Date()
	current, err := time.Parse("2006-01-02", "2021-01-01")
	if err != nil {
		return nil, err
	}

	for {
		date := current.Format("2006-01-02")
		dates = append(dates, date)

		currYear, currMonth, curDay := current.Date()
		if currYear == maxYear && currMonth == maxMonth && curDay == maxDay {
			break
		}

		current = current.AddDate(0, 0, 1)
	}

	return dates, nil
}
