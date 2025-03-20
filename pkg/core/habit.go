package core

type Habit struct {
	ID    int32
	Date  Date
	state HabitState
	Kind  HabitType
}

type HabitCategory struct {
	ID          int32
	Kind        HabitType
	Name        string
	Description string
}
