package sqlite

import (
	"database/sql"
	"strings"
	"time"

	"github.com/danielcosme/curious-ape/internal/core/database"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/go-sdk/errors"
	"github.com/jmoiron/sqlx"
)

type DaysDataSource struct {
	DB *sqlx.DB
}

func (ds *DaysDataSource) Create(d *entity.Day) error {
	d.Date = sanitizeDate(d.Date)
	query := `
		INSERT INTO "days" ("date") 
		VALUES (:date);
	`
	_, err := ds.DB.NamedExec(query, d)
	return err
}

func (ds *DaysDataSource) Update(date *entity.Day, joins ...entity.DayJoin) (*entity.Day, error) {
	query := `
		UPDATE "days"
		SET deep_work_minutes = :deep_work_minutes
		WHERE id = :id
    `
	_, err := ds.DB.NamedExec(query, date)
	if err != nil {
		return nil, catchErr(err)
	}
	return ds.Get(entity.DayFilter{IDs: []int{date.ID}}, joins...)
}

func (ds *DaysDataSource) Get(filter entity.DayFilter, joins ...entity.DayJoin) (*entity.Day, error) {
	day := new(entity.Day)
	q, args := dayFilter(filter).generate()
	if err := ds.DB.Get(day, q, args...); err != nil {
		return nil, catchErr(err)
	}
	return day, database.ExecuteDaysPipeline([]*entity.Day{day}, joins...)
}

func (ds *DaysDataSource) Find(filter entity.DayFilter, joins ...entity.DayJoin) ([]*entity.Day, error) {
	days := []*entity.Day{}
	q, args := dayFilter(filter).generate()
	if err := ds.DB.Select(&days, q, args...); err != nil {
		return nil, err
	}
	return days, database.ExecuteDaysPipeline(days, joins...)
}

func catchErr(err error) error {
	if err == nil {
		return nil
	}

	switch err.Error() {
	case sql.ErrNoRows.Error():
		return database.ErrNotFound
	default:
		if strings.Contains(err.Error(), "UNIQUE constraint failed:") {
			return database.ErrUniqueCheckFailed
		}
		return errors.NewFatal(err.Error())
	}
}

func dayFilter(f entity.DayFilter) *sqlQueryBuilder {
	b := newBuilder("days")

	if len(f.IDs) > 0 {
		b.AddFilter("id", intToInterface(f.IDs))
	}

	if len(f.Dates) > 0 {
		b.AddFilter("date", dateToInterface(f.Dates))
	}

	return b
}

func sanitizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
