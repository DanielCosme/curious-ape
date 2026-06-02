package core

import (
	"errors"
)

type Deadline struct {
	RepositoryCommon
	Title     string
	StartDate Date
	EndDate   Date
	Recurring bool
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
	}
	return nil
}
