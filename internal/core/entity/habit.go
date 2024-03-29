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

func CalculateHabitStatus(logs []*HabitLog) HabitStatus {
	status := HabitStatusNoInfo
	var override bool
	for _, log := range logs {
		if !log.IsAutomated {
			if log.Success {
				return HabitStatusDone
			}
			status = HabitStatusNotDone
			override = true
			continue
		}

		if !override {
			if log.Success {
				status = HabitStatusDone
			} else if status != HabitStatusDone {
				status = HabitStatusNotDone
			}
		}
	}
	return status
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
	ID          int        `db:"id"`
	HabitID     int        `db:"habit_id"`
	Success     bool       `db:"success"`
	Note        string     `db:"note"`
	Origin      DataSource `db:"origin"`
	IsAutomated bool       `db:"is_automated"`
}

type HabitType string

const (
	HabitTypeFood     HabitType = "food"
	HabitTypeWakeUp   HabitType = "wake_up"
	HabitTypeFitness  HabitType = "fitness"
	HabitTypeDeepWork HabitType = "deep_work"
	HabitTypeCustom   HabitType = "custom"
)

func (ht HabitType) Str() string {
	return string(ht)
}

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
	ID   []int
	Type []HabitType
	Code []string
}

type HabitLogFilter struct {
	ID      []int
	HabitID []int
	Origin  []DataSource
}

type HabitJoin func(hs []*Habit) error
