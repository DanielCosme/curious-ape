package core

type FitnessLogType string

const (
	FitnessLogTypeStrength FitnessLogType = "strength"
	FitnessLogTypeCardio   FitnessLogType = "cardio"
)

type FitnessLog struct {
	RepositoryCommon
	TimelineLog
	Date   Date
	Type   FitnessLogType
	Origin LogOrigin
	Raw    []byte
}
