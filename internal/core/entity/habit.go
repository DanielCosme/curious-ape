package entity

import (
	habitsC "github.com/danielcosme/curious-ape/internal/core/entity/constants/habit"
	"time"
)

type Habit struct {
	Entity
	Day     Day
	Time    time.Time
	Type    HabitType
	History []HabitHistoryEntry
	// generated
	State habitsC.State
}

type HabitType struct {
	Entity
	Name string
	Code habitsC.Code
}

type HabitHistoryEntry struct {
	Entity
	HabitID     ID
	Name        string // provider, client, automated? (fitbit, google, web, android, system, cron, event)
	Automated   bool   // manual or automated
	Success     bool   // yes or not
	Description string // by missing record, by waking up before 7, 5hours of work, because I can, etc.
}

type HabitQuery struct {
	*Query
	*DateQuery
	Code habitsC.Code
}
