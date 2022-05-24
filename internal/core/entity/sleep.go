package entity

type SleepRecord struct {
	// Repository
	ID int `db:"id"`

	TimeInBed int
	TimeAwake int
}
