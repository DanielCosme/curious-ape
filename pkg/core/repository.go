package core

type DayRepository interface {
	Create(Date) (Day, error)
	Get(DayParams) (Day, error)
	GetOrCreate(DayParams) (Day, error)
	Find(DayParams) ([]Day, error)
}

type HabitRepository interface {
	Get(HabitParams) (Habit, error)
	Upsert(UpsertHabitParams) (Habit, error)
}

type RepositoryCommon struct {
	ID uint
}
