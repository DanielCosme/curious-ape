package sync

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/danielcosme/curious-ape/internal/data"
	"github.com/danielcosme/curious-ape/internal/sync/fitbit"
	"github.com/danielcosme/curious-ape/internal/sync/toggl"
)

var (
	ErrTokenExpired = errors.New("token expired")
	ErrUnauthorized = errors.New("server needs to authorize again")
)

// type Collectors interface {
// 	LogsRange()
// 	DayLog()
// }

type Collector struct {
	Models *data.Models
	Sleep  *fitbit.SleepCollector
	Work   *toggl.WorkCollector
}

func NewCollectors(models *data.Models) *Collector {
	f := &fitbit.SleepCollector{
		Auth:  fitbit.FitbitAuth,
		Token: &models.Tokens,
		Scope: "sleep",
	}

	togg := &toggl.WorkCollector{
		Auth:  toggl.TogglAuth,
		Scope: "work",
	}

	return &Collector{
		Models: models,
		Sleep:  f,
		Work:   togg,
	}
}

func (co *Collector) GetTodayLog() error {
	today := time.Now()
	strTime := today.Format("2006-01-02")
	err := co.GetLog(strTime)
	if err != nil {
		return err
	}
	return nil
}

func (co *Collector) GetLog(date string) error {
	log.Println("Getting Record")
	jsonRecord, err := co.Sleep.DayLog(date)
	if err != nil {
		log.Println("ERR on GetLog", err.Error())
		return err
	}

	err = co.saveLog(date, jsonRecord)
	if err != nil {
		log.Println("ERR on GetLog", err.Error())
		return err
	}

	return nil
}

func (co *Collector) FromDayZero(limit time.Time) error {
	zero, _ := time.Parse("2006-01-02", fitbit.ZeroDay) // From here
	err := co.AllRangeLogs(zero, limit)
	if err != nil {
		return err
	}
	return nil
}

func (co *Collector) AllRangeLogs(zero, limit time.Time) error {
	maxYear, maxMonth, _ := limit.Date()

	first := zero
	for {
		last := first.AddDate(0, 1, -1)
		firstStr := first.Format("2006-01-02")
		lastStr := last.Format("2006-01-02")

		// Get logs for current month
		if err := co.monthRangeLog(firstStr, lastStr); err != nil {
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
func (co *Collector) monthRangeLog(first, last string) error {
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
func (co *Collector) saveLog(date string, jsonResponse []byte) error {
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

	err = co.SaveSleepHabit(&sleepRecord)
	if err != nil {
		log.Println("ERR", err.Error())
		return err
	}

	log.Println("All Good")
	return nil
}

func (co *Collector) SaveSleepHabit(sleepRecord *data.SleepRecord) error {
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

func (co *Collector) InitializeDayHabits() (err error) {
	types := []string{"sleep", "food", "fitness", "work"}
	h := data.Habit{
		State:  "no_info",
		Date:   time.Now().Format("2006-01-02"),
		Origin: "automated",
	}

	c := 0
	for _, v := range types {
		h.Type = v
		err = co.Models.Habits.Insert(&h)
		if err == nil {
			c++
		}
	}

	if err != nil {
		log.Println(c, "Habits Added,", err.Error())
		return err
	}

	log.Println(c, "CRON habits for today added successfully")
	return nil
}

func (co *Collector) BuildHabitsFromSleepRecords() (err error) {
	all, err := co.Models.SleepRecords.GetAll()
	if err != nil {
		return err
	}

	l := len(all)
	for i, v := range all {
		err := co.SaveSleepHabit(v)
		if err != nil {
			return err
		}

		if i+1 == l {
			log.Println("succesfully added", l, "habits")
		}
	}

	return nil
}
