package application

import "time"

type Habit struct {
	State HabitState
	Type  HabitType
	Date  time.Time
}
type HabitState string

const (
	HabitStateDone    HabitState = "done"
	HabitStateNotDone HabitState = "not-done"
	HabitStateNoInfo  HabitState = "no-info"
)

func StateValid(s HabitState) bool {
	switch s {
	case HabitStateDone, HabitStateNotDone, HabitStateNoInfo:
		return true
	}
	return false
}

type HabitType string

const (
	HabitTypeWake    HabitType = "wake"
	HabitTypeWorkout HabitType = "workout"
	HabitTypeWork    HabitType = "work"
	HabitTypeEat     HabitType = "eat"
)

func HabitTypeValid(s HabitType) bool {
	switch s {
	case HabitTypeWake, HabitTypeWorkout, HabitTypeWork, HabitTypeEat:
		return true
	}
	return false
}

func (app *Application) HabitUpsert(date time.Time, state HabitState, habitType HabitType) (*Habit, error) {
	return &Habit{
		State: state,
		Type:  habitType,
		Date:  date,
	}, nil
}
