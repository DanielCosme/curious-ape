package core

type Day struct {
	RepositoryCommon
	Date   Date
	Habits []Habit
}

func (d *Day) IsZero() bool {
	return d.Date.Time().IsZero()
}

type DayParams struct {
	ID    uint
	Date  Date
	Dates DateSlice
	Lite  bool
}

type FitnessLog struct{}
type DeepWorkLog struct{}
type SleepLog struct{}
