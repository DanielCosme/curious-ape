package core

import (
	"fmt"
	"time"
)

type Date struct {
	time time.Time
}

type DateSlice []Date

func NewDate(t time.Time) Date {
	return Date{
		time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local),
	}
}

func (d Date) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%q", d.String())
	return []byte(s), nil
}

func NewDateToday() Date {
	t := time.Now().Local()
	return Date{
		time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local),
	}
}

func (d Date) IsEqual(t time.Time) bool {
	return d.Time().Equal(NewDate(t).Time())
}

func (d Date) String() string {
	return d.time.Format(ISO8601)
}

func (d Date) Time() time.Time {
	return d.time
}

func (d Date) FirstDayOfTheMonth() Date {
	return NewDate(time.Date(d.Time().Year(), d.Time().Month(), 1, 0, 0, 0, 0, d.time.Location()))
}

func (d Date) RangeMonth() DateSlice {
	var dates []Date
	beginning := d.FirstDayOfTheMonth().Time()

	for beginning.Before(d.time) {
		dates = append(dates, NewDate(beginning))
		beginning = beginning.AddDate(0, 0, 1)
	}

	dates = append(dates, d)
	return dates

}

func (d Date) ToEndOfDay() time.Time {
	t := d.Time()
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
}

func (d Date) ToBeginningOfDay() time.Time {
	t := d.Time()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func (ds DateSlice) ToTimeSlice() []time.Time {
	res := make([]time.Time, len(ds))
	for idx, d := range ds {
		res[idx] = d.Time()
	}
	return res
}

func DateFromISO8601(s string) (Date, error) {
	t, err := time.Parse(ISO8601, s)
	if err != nil {
		return Date{}, err
	}
	return NewDate(t), nil
}
