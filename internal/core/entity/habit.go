package entity

const HabitCodeDefault = "default"

type Habit struct {
	// Repository
	ID         int `db:"id"`
	DayID      int `db:"day_id"`
	CategoryID int `db:"habit_category_id"`
	// Core
	Status HabitStatus `db:"status"`
	// Generated
	Day      *Day
	Category *HabitCategory
	Logs     []*HabitLog
}

type HabitCategory struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Type        HabitType `db:"type"`
	Code        string    `db:"code"`
	Description string    `db:"description"`
	Color       string    `db:"color"`
}

type HabitLog struct {
	ID          int         `db:"id"`
	HabitID     int         `db:"habit_id"`
	Note        string      `db:"note"`
	Origin      HabitOrigin `db:"origin"`
	IsAutomated bool        `db:"is_automated"`
	Success     bool        `db:"success"`
}

type HabitType string

const (
	HabitTypeFood     HabitType = "food"
	HabitTypeWakeUp   HabitType = "wake-up"
	HabitTypeFitness  HabitType = "fitness"
	HabitTypeDeepWork HabitType = "deep-work"
	HabitTypeCustom   HabitType = "custom"
)

func IsValidHabitCategoryType(habitType HabitType) bool {
	switch habitType {
	case HabitTypeFood, HabitTypeCustom, HabitTypeFitness, HabitTypeWakeUp, HabitTypeDeepWork:
		return true
	}
	return false
}

type HabitOrigin string

const (
	HabitOriginClient    HabitOrigin = "client"   // android/web/cli
	HabitOriginProvider  HabitOrigin = "provider" // fitbit, google etc
	HabitOriginWebSystem HabitOrigin = "system"   // internal from manual entries?
	HabitOriginUnknown   HabitOrigin = "unknown"  // ??
)

type HabitStatus string

const (
	HabitStatusDone    HabitStatus = "done"
	HabitStatusNotDone HabitStatus = "not-done"
	HabitStatusNoInfo  HabitStatus = "no-info"
)

type HabitFilter struct {
	ID         []int
	DayID      []int
	CategoryID []int
}

type HabitCategoryFilter struct {
	ID []int
}

type HabitLogFilter struct {
	ID      []int
	HabitID []int
	Origin  []HabitOrigin
}

type HabitJoin func(hs []*Habit) error
