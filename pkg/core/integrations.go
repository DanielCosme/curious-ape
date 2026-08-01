package core

type WorkLogIntegration interface {
	GetDayEntries(Date) ([]DeepWorkLog, error)
}
