package core

import "time"

type DeepWorkLog struct {
	RepositoryCommon
	TimelineLog
	Date   Date
	Origin LogOrigin
	Raw    []byte
}

type DeepWorkLogParams struct {
	ID        int64
	DayID     int64
	StartTime time.Time
}
