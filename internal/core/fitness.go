package core

type FitnessRecord struct {
	ID                  int    `json:"id"`
	Date                string `json:"date"`
	StartInMilliseconds int    `json:"startInMilliseconds"` // timestamp
	EndInMilliseconds   int    `json:"endInMilliseconds"`
	Provider            string `json:"provider"`
}

type FitnessModel interface {
	Insert(f *FitnessRecord) error
	GetAll() ([]*FitnessRecord, error)
	Get(date string) (*FitnessRecord, error)
}
