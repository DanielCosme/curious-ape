package core

type DayRepository interface {
	Create(Date) (Day, error)
	Get(DayParams) (Day, error)
	Find(DayParams) ([]Day, error)
}

type HabitRepository interface {
	Get(HabitParams) (Habit, error)
	Upsert(Habit) (Habit, error)
}

type SleepLogRepository interface {
	Upsert(SleepLog) (SleepLog, error)
}

type FitnessLogRepository interface {
	Upsert(FitnessLog) (FitnessLog, error)
}

type RepositoryCommon struct {
	ID uint
}
