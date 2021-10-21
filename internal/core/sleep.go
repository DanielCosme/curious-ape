package core

type SleepRecord struct {
	ID            int    `json:"id"`
	Date          string `json:"dateOfSleep"`
	Duration      int    `json:"duration"`
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
	MinutesAsleep int    `json:"minutesAsleep"`
	MinutesAwake  int    `json:"minutesAwake"`
	MinutesInBed  int    `json:"timeInBed"`
	Provider      string `json:"-"`
	RawJson       []byte `json:"-"`
}

type SleepRecordModel interface {
	Get(date string) (*SleepRecord, error)
	GetAll() ([]*SleepRecord, error)
	Insert(data *SleepRecord) error
}
