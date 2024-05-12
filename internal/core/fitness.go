package core

import "time"

type FitnessLog struct {
	ID        int        `db:"id"`
	DayID     int        `db:"day_id"`
	Title     string     `db:"title"`
	Type      string     `db:"type"`
	Date      time.Time  `db:"date"`
	StartTime time.Time  `db:"start_time"`
	EndTime   time.Time  `db:"end_time"`
	Origin    DataSource `db:"origin"`
	Raw       string     `db:"raw"`
	Note      string     `db:"note"`
	// Generated
	Day *Day
}
