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

func (sl *SleepLog) ToHabitLogWakeUp() (res NewHabitParams) {
	res.Date = sl.Date
	res.HabitType = HabitTypeWakeUp
	res.Origin = OriginLogSleep
	res.Automated = true
	res.Detail = sl.EndTime.Format(Time)
	wakeUpTime := time.Date(sl.EndTime.Year(), sl.EndTime.Month(), sl.EndTime.Day(), 6, 0, 0, 0, sl.EndTime.Location())
	if sl.EndTime.Before(wakeUpTime) {
		res.Success = true
	}
	return res
}
