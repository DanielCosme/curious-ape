package core

import "time"

type TimelineType uint

const (
	TimelineTypeSleep TimelineType = iota + 1
	TimelineTypeFitness
	TimelineTypeDeepWork
)

type Timeliner interface {
	Timeline() TimelineLog
}

type TimelineLog struct {
	Title     string
	StartTime time.Time
	EndTime   time.Time
	Note      string
	Type      TimelineType
}
