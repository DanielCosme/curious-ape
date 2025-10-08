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

type DaySliceSortDESC []Day

// NOTE: Sorts DESC
func (ds DaySliceSortDESC) Less(i, j int) bool { return ds[i].Date.Time().After(ds[j].Date.Time()) }
func (ds DaySliceSortDESC) Swap(i, j int)      { ds[i], ds[j] = ds[j], ds[i] }
func (ds DaySliceSortDESC) Len() int           { return len(ds) }
