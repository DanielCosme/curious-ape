package core

type HabitState string

const (
	HabitStateDone    HabitState = "done"
	HabitStateNotDone HabitState = "not_done"
	HabitStateNoInfo  HabitState = "no_info"
)

type HabitType string

const (
	HabitTypeWakeUp     HabitType = "wake_up"
	HabitTypeFitness    HabitType = "fitness"
	HabitTypeDeepWork   HabitType = "deep_work"
	HabitTypeEatHealthy HabitType = "food"
)

func HabitTypeFromString(s string) HabitType {
	switch s {
	case "food":
		return HabitTypeEatHealthy
	case "wake_up":
		return HabitTypeWakeUp
	case "fitness":
		return HabitTypeFitness
	case "deep_work":
		return HabitTypeDeepWork
	}
	return "undefined"
}

type Habit struct {
	RepositoryCommon
	Date      Date
	State     HabitState
	Type      HabitType
	Automated bool
}

type HabitCategory struct {
	RepositoryCommon
	Name        string
	Kind        HabitType
	Description string
}

type HabitParams struct {
	ID int
}

type UpsertHabitParams struct {
	Date      Date
	Type      HabitType
	State     HabitState
	Automated bool
}

type HabitCategoryParams struct {
	ID   int
	Kind HabitType
}
