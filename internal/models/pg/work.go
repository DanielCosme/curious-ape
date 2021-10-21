package pg

import (
	"database/sql"
	"github.com/danielcosme/curious-ape/internal/core"
)

type WorkModel struct {
	DB *sql.DB
}

func (wm *WorkModel) Insert(wr *core.WorkRecord) error {
	stm := `INSERT INTO work_records ("date", "grand_total", raw_json, provider)
			VALUES($1, $2, $3, $4)`
	_, err := wm.DB.Exec(stm, wr.Date, wr.Total, wr.RawJson, wr.Provider)
	if err != nil {
		return err
	}
	return nil
}

func (wm *WorkModel) Get(date string) (*core.WorkRecord, error) {
	r := &core.WorkRecord{}
	stm := `SELECT id, "date", grand_total, provider FROM work_records WHERE "date" = $1`
	row := wm.DB.QueryRow(stm, date)

	err := row.Scan(&r.ID, &r.Date, &r.Total, &r.Provider)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (wm *WorkModel) GetAll() ([]*core.WorkRecord, error) {
	records := []*core.WorkRecord{}
	stm := `Select id, "date", grand_total, provider from work_records`
	rows, err := wm.DB.Query(stm)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		r := &core.WorkRecord{}
		err := rows.Scan(&r.ID, &r.Date, &r.Total, &r.Provider)
		if err != nil {
			return nil, err
		}
		records = append(records, r)
	}

	return records, nil
}
