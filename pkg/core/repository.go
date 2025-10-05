package core

type DayRepository interface {
	Create(Date) (Day, error)
	Get(DayParams) (Day, error)
	Find(DayParams) ([]Day, error)
}

type HabitRepository interface {
	Get(HabitParams) (Habit, error)
	Upsert(HabitUpsertParams) (Habit, error)
}

type SleepLogRepository interface {
	Upsert(SleepLogUpsertParams) (SleepLog, error)
}

type RepositoryCommon struct {
	ID uint
}
