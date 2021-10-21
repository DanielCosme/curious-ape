package core

type WorkRecord struct {
	ID       int    `json:"id"`
	Date     string `json:"date"`
	Total    int    `json:"total_grand"`
	RawJson  string `json:"rawJson"`
	Provider string `json:"provider"`
}

type WorkModel interface {
	Insert(wr *WorkRecord) error
	Get(date string) (*WorkRecord, error)
	GetAll() ([]*WorkRecord, error)
}
