package fitbit

import (
	"time"

	"git.danicos.dev/daniel/curious-ape/pkg/core"
)

func ToDuration(i int) time.Duration {
	return time.Duration(i) * time.Minute
}

func ParseDate(s string) time.Time {
	// yyyy-mm-dd
	t, _ := time.Parse("2006-01-02", s)
	return t
}

func ParseTime(s string) time.Time {
	// 2022-06-02T05:18:30.000
	t, _ := time.Parse("2006-01-02T15:04:05.999", s)
	return core.TimeUTC(t)
}
