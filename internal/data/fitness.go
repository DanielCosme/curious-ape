package data

import (
	"database/sql"
	"errors"
)

type FitnessRecord struct {
	ID                 int    `json:"id"`
	Date               string `json:"date"`
	StartInMiliseconds int    `json:"startInMilisseconds"` // timestamp
	EndInMiliseconds   int    `json:"endInMiliseconds"`
	Provider           string `json:"provider"`
}

type FitnessModel struct {
	DB *sql.DB
}

func (fm *FitnessModel) Insert(f *FitnessRecord) error {
	stm := `INSERT INTO fitness_records
			("date", start_in_miliseconds, end_in_miliseconds, provider)
			VALUES ($1, $2, $3, $4)`
	_, err := fm.DB.Exec(stm, f.Date, f.StartInMiliseconds, f.EndInMiliseconds, f.Provider)
	return err
}

func (fm *FitnessModel) Get(date string) (*FitnessRecord, error) {
	r := &FitnessRecord{}
	stm := `SELECT id, "date", start_in_miliseconds, end_in_miliseconds, provider
			FROM fitness_records WHERE "date" = $1`
	row := fm.DB.QueryRow(stm, date)
	err := row.Scan(&r.ID, &r.Date, &r.StartInMiliseconds, &r.EndInMiliseconds, &r.Provider)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return r, nil
}

func (fm *FitnessModel) GetAll() ([]*FitnessRecord, error) {
	records := []*FitnessRecord{}
	stm := `SELECT * FROM fitness_records`
	rows, err := fm.DB.Query(stm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := &FitnessRecord{}
		err := rows.Scan(
			&r.ID,
			&r.Date,
			&r.StartInMiliseconds,
			&r.EndInMiliseconds,
			&r.Provider,
		)

		if err != nil {
			return nil, err
		}

		records = append(records, r)
	}

	return records, nil
}
