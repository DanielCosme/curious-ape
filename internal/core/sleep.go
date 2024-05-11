package core

import "time"

type SleepLog struct {
	ID          int        `db:"id"`
	DayID       int        `db:"day_id"`
	Date        time.Time  `db:"date"`
	StartTime   time.Time  `db:"start_time"`
	EndTime     time.Time  `db:"end_time"`
	IsMainSleep bool       `db:"is_main_sleep"`
	IsAutomated bool       `db:"is_automated"`
	Origin      DataSource `db:"origin"`
	Raw         string     `db:"raw"`
	// Minutes
	MinutesInBed  time.Duration `db:"total_time_in_bed"`
	MinutesAsleep time.Duration `db:"minutes_asleep"`
	MinutesRem    time.Duration `db:"minutes_rem"`
	MinutesDeep   time.Duration `db:"minutes_deep"`
	MinutesLight  time.Duration `db:"minutes_light"`
	MinutesAwake  time.Duration `db:"minutes_awake"`
	// Generated
	Day *Day
}
