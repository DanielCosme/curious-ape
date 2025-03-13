package core

type Habit struct {
	ID    int32
	Date  Date
	state HabitState
	Kind  HabitKind
}

type HabitCategory struct {
	ID          int32
	Kind        HabitKind
	Name        string
	Description string
}
