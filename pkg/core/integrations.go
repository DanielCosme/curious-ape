package core

type WorkLogIntegration interface {
	GetDayEntries(Date) ([]DeepWorkLog, error)
}

type FitnessIntegration interface {
	GetFitnessLogs(Date) ([]FitnessLog, error)
}

type SleepIntegration interface {
	GetSleepLogs(Date) ([]FitnessLog, error)
}
