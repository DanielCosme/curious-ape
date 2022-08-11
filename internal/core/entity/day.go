package entity

import "time"

const ISO8601 = "2006-01-02"
const HumanDate = "Monday, 02 January 2006"
const Timestamp = "15:04:05"
const Time = "15:04"

type Day struct {
	ID              int       `db:"id"`
	Date            time.Time `db:"date"`
	DeepWorkMinutes int       `db:"deep_work_minutes"`
	// generated
	Habits      []*Habit
	SleepLogs   []*SleepLog
	FitnessLogs []*FitnessLog
	Tags        []*Tag
	// Days will have tags
	// Should the tags by their own entity, or just an array of strings?
	// e.g [Sick, BadDay, Sunrise, CA, CO]
	// Country Tags??
	// Mood Tags??
	// Weather Tags??
	// Health Tags?? -> Covid, Covid-Recovery
	// Weather && Location on that day
	// Location -> Where did I wake up && went to sleep this day?
	// Tags should be dynamic and be able to be created by me ay any time.
	// Habit Tags?
}

func (d *Day) FormatDate() string {
	return d.Date.Format(ISO8601)
}

type Tag struct{}

type DayJoin func([]*Day) error

type DayFilter struct {
	IDs   []int
	Dates []time.Time
}

func ParseDate(d string) (time.Time, error) {
	var t time.Time
	t, err := time.Parse(ISO8601, d)
	if err != nil {
		return t, err
	}
	return NormalizeDate(t, time.UTC), nil
}

func ParseTime(t string) (time.Time, error) {
	return time.Parse(Time, t)
}

func FormatDate(t time.Time) string {
	return t.Format(ISO8601)
}

func NormalizeDate(t time.Time, loc *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
}
