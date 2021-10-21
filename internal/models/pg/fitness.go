package pg

import (
	"database/sql"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/errors"
)

type FitnessModel struct {
	DB *sql.DB
}

func (fm *FitnessModel) Insert(f *core.FitnessRecord) error {
	stm := `INSERT INTO fitness_records
			("date", start_in_miliseconds, end_in_miliseconds, provider)
			VALUES ($1, $2, $3, $4)`
	_, err := fm.DB.Exec(stm, f.Date, f.StartInMilliseconds, f.EndInMilliseconds, f.Provider)
	return err
}

func (fm *FitnessModel) Get(date string) (*core.FitnessRecord, error) {
	r := &core.FitnessRecord{}
	stm := `SELECT id, "date", start_in_miliseconds, end_in_miliseconds, provider
			FROM fitness_records WHERE "date" = $1`
	row := fm.DB.QueryRow(stm, date)
	err := row.Scan(&r.ID, &r.Date, &r.StartInMilliseconds, &r.EndInMilliseconds, &r.Provider)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return r, nil
}

func (fm *FitnessModel) GetAll() ([]*core.FitnessRecord, error) {
	records := []*core.FitnessRecord{}
	stm := `SELECT * FROM fitness_records`
	rows, err := fm.DB.Query(stm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := &core.FitnessRecord{}
		err := rows.Scan(
			&r.ID,
			&r.Date,
			&r.StartInMilliseconds,
			&r.EndInMilliseconds,
			&r.Provider,
		)

		if err != nil {
			return nil, err
		}

		records = append(records, r)
	}

	return records, nil
}