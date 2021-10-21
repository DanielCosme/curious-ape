package pg

import (
	"database/sql"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/errors"
)

type SleepRecordModel struct {
	DB *sql.DB
}

func (sr *SleepRecordModel) Get(date string) (*core.SleepRecord, error) {
	r := &core.SleepRecord{}
	stm := `SELECT id, date, duration, start_time, end_time, minutes_asleep, minutes_awake, minutes_in_bed FROM sleep_records WHERE "date" = $1`
	row := sr.DB.QueryRow(stm, date)
	err := row.Scan(
		&r.ID,
		&r.Date,
		&r.Duration,
		&r.StartTime,
		&r.EndTime,
		&r.MinutesAsleep,
		&r.MinutesAwake,
		&r.MinutesInBed,
	)

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

func (sr *SleepRecordModel) GetAll() ([]*core.SleepRecord, error) {
	records := []*core.SleepRecord{}
	stm := `SELECT id, date, duration, start_time, end_time, minutes_asleep, minutes_awake, minutes_in_bed, provider FROM sleep_records`
	rows, err := sr.DB.Query(stm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := &core.SleepRecord{}
		err := rows.Scan(
			&r.ID,
			&r.Date,
			&r.Duration,
			&r.StartTime,
			&r.EndTime,
			&r.MinutesAsleep,
			&r.MinutesAwake,
			&r.MinutesInBed,
			&r.Provider,
		)

		if err != nil {
			return nil, err
		}

		records = append(records, r)
	}

	return records, nil
}

func (sr *SleepRecordModel) Insert(data *core.SleepRecord) error {
	stm := `INSERT INTO sleep_records (date, duration, start_time, end_time, minutes_asleep,
									minutes_awake, minutes_in_bed, provider, raw_json)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	args := []interface{}{
		data.Date,
		data.Duration,
		data.StartTime,
		data.EndTime,
		data.MinutesAsleep,
		data.MinutesAwake,
		data.MinutesInBed,
		data.Provider,
		data.RawJson,
	}
	_, err := sr.DB.Exec(stm, args...)
	if err != nil {
		return err
	}

	return nil
}
