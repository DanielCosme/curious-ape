package entity

type Habit struct {
	Entity
	Day   Day
	State HabitState `db:"state"`
	Type  *HabitType
	// History []*HabitHistoryEntry
}

type HabitType struct {
	Entity
	Name string
	Code HabitCode
}

type HabitHistoryEntry struct {
	Entity
	HabitID     string
	Name        string // provider, client, automated? (fitbit, google, web, android, system, cron, event)
	Automated   bool   // manual or automated
	Success     bool   // yes or not
	Description string // by missing record, by waking up before 7, 5hours of work, because I can, etc.
}

type HabitQuery struct {
	*Query
	*DateQuery
	Code HabitCode
}

type HabitState string

const (
	Done    HabitState = "done"
	NoInfo  HabitState = "no-info"
	NotDone HabitState = "not-done"
)

type HabitCode string

const (
	Food     HabitCode = "food"
	DeepWork HabitCode = "deep-work"
	Fitness  HabitCode = "fitness"
	WakeUp   HabitCode = "wake-up"
	Custom   HabitCode = "custom"
)
