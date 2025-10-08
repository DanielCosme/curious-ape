package core

import (
	"time"
)

type SleepLog struct {
	RepositoryCommon
	TimelineLog
	Date        Date
	IsMainSleep bool
	TimeAsleep  time.Duration
	TimeInBed   time.Duration
	Origin      LogOrigin
	Raw         []byte
}
