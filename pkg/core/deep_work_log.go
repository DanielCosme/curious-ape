package core

type DeepWorkLog struct {
	RepositoryCommon
	TimelineLog
	Date   Date
	Origin LogOrigin
	Raw    []byte
}
