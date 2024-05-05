package dates

import "time"

func ToUTC(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
}

func ToLocal(t time.Time, loc *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
}

func Ptr(t time.Time) *time.Time {
	return &t
}

func BeginningOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 0, 0, t.Location())
}

func StartAndEndOfweek(t time.Time) (time.Time, time.Time) {
	return BeginningOfDay(t), EndOfWeek(t)
}

func BeginningOfWeek(t time.Time) time.Time {
	for t.Weekday() != time.Monday {
		t = t.AddDate(0, 0, -1)
	}
	return t
}

func EndOfWeek(t time.Time) time.Time {
	for t.Weekday() != time.Sunday {
		t = t.AddDate(0, 0, 1)
	}
	return t
}

func StartAndEndOfTheMonth(t time.Time) (time.Time, time.Time) {
	return StartOfTheMonth(t), EndOfTheMonth(t)
}

func StartOfTheMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

func EndOfTheMonth(t time.Time) time.Time {
	t = t.AddDate(0, 1, 0)
	t = time.Date(t.Year(), t.Month(), 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	return t.AddDate(0, 0, -1)
}
