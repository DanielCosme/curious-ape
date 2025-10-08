package core

type DayRelations int

const (
	DayRelationHabits DayRelations = iota
	DayRelationFitnessLogs
	DayRelationDeepWorkLogs
	DayRelationSleepLogs
)

type Day struct {
	RepositoryCommon
	Date         Date
	Habits       DayHabits
	SleepLogs    []SleepLog
	FitnessLogs  []FitnessLog
	DeepWorkLogs []DeepWorkLog
}

type DayHabits struct {
	Hs       []Habit
	Sleep    Habit
	Fitness  Habit
	DeepWork Habit
	Eat      Habit
}

func (d *Day) IsZero() bool {
	return d.Date.Time().IsZero()
}

type DayParams struct {
	ID    uint
	Date  Date
	Dates DateSlice
	Order OrderParam
}

type OrderParam int

const (
	ASC OrderParam = iota
	DESC
)

func DayRelationsAll() []DayRelations {
	return []DayRelations{
		DayRelationHabits,
		DayRelationFitnessLogs,
		DayRelationDeepWorkLogs,
		DayRelationSleepLogs,
	}
}
