package data

import (
	"database/sql"
)

type WorkRecord struct {
	ID       int    `json:"id"`
	Date     string `json:"date"`
	Total    int    `json:"total_grand"`
	RawJson  string `json:"rawJson"`
	Provider string `json:"provider"`
}

type WorkModel struct {
	DB *sql.DB
}

func (wm *WorkModel) Insert(wr *WorkRecord) error {
	stm := `INSERT INTO work_records ("date", "grand_total", raw_json, provider)
			VALUES($1, $2, $3, $4)`
	_, err := wm.DB.Exec(stm, wr.Date, wr.Total, wr.RawJson, wr.Provider)
	if err != nil {
		return err
	}
	return nil
}

func (wm *WorkModel) Get(date string) (*WorkRecord, error) {
	r := &WorkRecord{}
	stm := `SELECT id, "date", grand_total, provider FROM work_records WHERE "date" = $1`
	row := wm.DB.QueryRow(stm, date)

	err := row.Scan(&r.ID, &r.Date, &r.Total, &r.Provider)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (wm *WorkModel) GetAll() ([]*WorkRecord, error) {
	records := []*WorkRecord{}
	stm := `Select id, "date", grand_total, provider from work_records`
	rows, err := wm.DB.Query(stm)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		r := &WorkRecord{}
		err := rows.Scan(&r.ID, &r.Date, &r.Total, &r.Provider)
		if err != nil {
			return nil, err
		}
		records = append(records, r)
	}

	return records, nil
}
