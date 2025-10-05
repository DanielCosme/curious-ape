package core

import "time"

type SleepLog struct {
	TimelineLog
	Date        Date
	IsMainSleep bool
	TimeAsleep  time.Duration
	TimeInBed   time.Duration
}

type SleepLogUpsertParams struct {
	TimelineLog
	Date        Date
	IsMainSleep bool
	TimeAsleep  time.Duration
	TimeInBed   time.Duration
	Origin      LogOrigin
	Raw         string
}
