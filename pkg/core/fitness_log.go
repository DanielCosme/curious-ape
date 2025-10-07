package core

type FitnessLogType string

const (
	FitnessLogTypeStrength FitnessLogType = "strength"
	FitnessLogTypeCardio   FitnessLogType = "cardio"
)

type FitnessLog struct {
	RepositoryCommon
	TimelineLog
	Date        Date
	FitnessType FitnessLogType
	Origin      LogOrigin
	Raw         []byte
}
