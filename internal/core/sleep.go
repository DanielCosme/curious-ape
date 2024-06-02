package core

import (
	"time"
)

type SleepLog struct {
	ID    int32
	DayID int32

	Date        Date
	StartTime   time.Time
	EndTime     time.Time
	IsMainSleep bool
	IsAutomated bool
	Origin      Integration
	Raw         []byte
	// Minutes
	MinutesInBed  time.Duration
	MinutesAsleep time.Duration
	MinutesRem    time.Duration
	MinutesDeep   time.Duration
	MinutesLight  time.Duration
	MinutesAwake  time.Duration
	// Generated
	Day *Day
}

func NewSleepLog(day Day) SleepLog {
	sl := SleepLog{
		DayID: day.ID,
		Date:  day.Date,
		Day:   &day,
	}
	return sl
}
