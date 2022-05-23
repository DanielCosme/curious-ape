package entity

const HabitCodeDefault = "default"

type Habit struct {
	// Repository
	ID         int `db:"id"`
	DayID      int `db:"day_id"`
	CategoryID int `db:"habit_category_id"`
	// Core
	IsAutomated bool        `db:"is_automated"`
	Success     bool        `db:"success"`
	Note        string      `db:"note"`
	Origin      HabitOrigin `db:"origin"`
	// Generated
	Category *HabitCategory `db:"habit_categories"`
	Day      *Day           `db:"day"`
}

type HabitCategory struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Type        HabitType `db:"type"`
	Code        string    `db:"code"`
	Description string    `db:"description"`
	Color       string    `db:"color"`
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
	HabitOriginClient    HabitOrigin = "client"
	HabitOriginProvider  HabitOrigin = "provider"
	HabitOriginWebSystem HabitOrigin = "system"
	HabitOriginUnknown   HabitOrigin = "unknown"
)

type HabitFilter struct {
	ID           []int
	CategoryIDs  []int
}

type HabitJoin func(hs []*Habit) error
