package core

import "time"

type Deadline struct {
	RepositoryCommon
	Title     string
	StartTime time.Time
	EndTime   time.Time
	Note      string
}
