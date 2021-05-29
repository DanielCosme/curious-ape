package data

import (
	"database/sql"
	"errors"
)

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

type SleepRecordModel struct {
	DB *sql.DB
}

func (sr *SleepRecordModel) Get(date string) (*SleepRecord, error) {
	r := &SleepRecord{}
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
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return r, nil
}

func (sr *SleepRecordModel) GetAll() ([]*SleepRecord, error) {
	records := []*SleepRecord{}
	stm := `SELECT id, date, duration, start_time, end_time, minutes_asleep, minutes_awake, minutes_in_bed FROM sleep_records`
	rows, err := sr.DB.Query(stm)
	if err != nil {
		return records, err
	}
	defer rows.Close()

	for rows.Next() {
		r := &SleepRecord{}
		err := rows.Scan(
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
			return records, err
		}

		records = append(records, r)
	}

	return records, err
}

func (sr *SleepRecordModel) Insert(data SleepRecord) error {
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
