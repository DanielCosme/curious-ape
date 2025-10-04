package core

import "time"

type SleepLog struct {
	TimelineLog
	Date        Date
	IsMainSleep bool
	TimeAsleep  time.Duration
	TimeInBed   time.Duration
}
