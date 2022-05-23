package entity

import "time"

const ISO8601 = "2006-01-02"

func ParseDate(d string) (time.Time, error) {
	var t time.Time
	t, err := time.Parse(ISO8601, d)
	if err != nil {
		return t, err
	}
	return NormalizeDate(t, time.UTC), nil
}

func FormatDate(t time.Time) string {
	return t.Format(ISO8601)
}

func NormalizeDate(t time.Time, loc *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
}

type Day struct {
	ID   int       `db:"id"`
	Date time.Time `db:"date"`
	// generated
	Habits []*Habit
}

type DayJoin func([]*Day) error

type DayFilter struct {
	IDs  []int
	Date []time.Time
}
