package core

import "time"

type FitnessLog struct {
	ID        int            `db:"id"`
	DayID     int            `db:"day_id"`
	Title     string         `db:"title"`
	Type      FitnessLogType `db:"type"`
	Date      time.Time      `db:"date"`
	StartTime time.Time      `db:"start_time"`
	EndTime   time.Time      `db:"end_time"`
	Origin    DataSource     `db:"origin"`
	Raw       string         `db:"raw"`
	Note      string         `db:"note"`
	// Generated
	Day *Day
}

type FitnessLogType string

const (
	StrengthTraining FitnessLogType = "strength_training"
	Steps            FitnessLogType = "steps"
	Cycling          FitnessLogType = "cycling"
	Run              FitnessLogType = "run"
)

func (flt FitnessLogType) Str() string {
	return string(flt)
}

type FitnessLogFilter struct {
	ID    []int
	DayID []int
	Date  []time.Time
}

type FitnessLogJoin func(hs []*FitnessLog) error
