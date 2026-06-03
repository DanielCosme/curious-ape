package core

import (
	"errors"
	"math"
)

type Deadline struct {
	RepositoryCommon
	Title     string
	StartDate Date
	EndDate   Date
	Recurring bool
	DaysLeft  int
}

func (d *Deadline) Validate() error {
	if d.Title == "" {
		return errors.New("title is empty")
	}
	if d.StartDate.Time().IsZero() {
		return errors.New("start time is empty")
	}
	if d.EndDate.Time().IsZero() {
		return errors.New("end time is empty")
	} else if d.EndDate.Time().Before(d.StartDate.Time()) {
		return errors.New("end time cannot be before start time")
	}
	return nil
}

// DaysLeft calculates days remaining from start until end time.
func DaysLeft(start, end Date) int {
	if start.time.After(end.time) {
		return 0 // Deadline has passed
	}
	remainingDuration := end.time.Sub(start.time)
	return int(math.Ceil(remainingDuration.Hours() / 24))
}
