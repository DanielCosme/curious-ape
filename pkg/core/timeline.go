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

type IntegrationStatus string

const IntegrationStatusConnected IntegrationStatus = "connected"
const IntegrationStatusUnkown IntegrationStatus = "unknown"
const IntegrationStatusDisconnected IntegrationStatus = "disconnected"
const IntegrationStatusNotImplemented IntegrationStatus = "not-implemented"

type IntegrationInfo struct {
	Name    string
	Status  IntegrationStatus
	Info    []string
	AuthURL string
}
