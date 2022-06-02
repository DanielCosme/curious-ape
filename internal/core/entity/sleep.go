package entity

import "time"

type SleepLog struct {
	// Repository
	ID          int            `db:"id"`
	DayID       int            `db:"day_id"`
	StartTime   time.Time      `db:"start_time"`
	EndTime     time.Time      `db:"end_time"`
	IsMainSleep bool           `db:"is_main_sleep"`
	IsAutomated bool           `db:"is_automated"`
	Origin      SleepLogOrigin `db:"origin"`
	Raw         string         `db:"raw"`
	// Minutes
	TotalTimeInBed int `db:"total_time_in_bed"`
	MinutesAsleep  int `db:"minutes_asleep"`
	MinutesRem     int `db:"minutes_rem"`
	MinutesDeep    int `db:"minutes_deep"`
	MinutesLight   int `db:"minutes_light"`
	MinutesAwake   int `db:"minutes_awake"`
	// Generated
	Day *Day
}

type SleepLogOrigin string

const (
	SleepLogOriginFitbit SleepLogOrigin = "fitbit"
	SleepLogOriginManual SleepLogOrigin = "manual"
)

type SleepLogFilter struct {
	ID    []int
	DayID []int
}

type SleepLogJoin func(hs []*SleepLog) error
