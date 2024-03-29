package sqlite

import (
	"database/sql"
	"fmt"
	logape "github.com/danielcosme/go-sdk/log"
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
	d.Date = normalizeDate(d.Date)
	query := `
		INSERT INTO "days" ("date") 
		VALUES (:date);
	`
	res, err := ds.DB.NamedExec(query, d)
	d.ID = lastInsertID(res)
	return catchErr("create day", err)
}

func (ds *DaysDataSource) Update(date *entity.Day, joins ...entity.DayJoin) (*entity.Day, error) {
	query := `
		UPDATE "days"
		SET deep_work_minutes = :deep_work_minutes
		WHERE id = :id
    `
	_, err := ds.DB.NamedExec(query, date)
	if err != nil {
		return nil, catchErr("update day", err)
	}
	return ds.Get(entity.DayFilter{IDs: []int{date.ID}}, joins...)
}

func (ds *DaysDataSource) Get(filter entity.DayFilter, joins ...entity.DayJoin) (*entity.Day, error) {
	day := new(entity.Day)
	q, args := dayFilter(filter).generate()
	if err := ds.DB.Get(day, q, args...); err != nil {
		return nil, catchErr("get day", err)
	}
	return day, database.ExecuteDaysPipeline([]*entity.Day{day}, joins...)
}

func (ds *DaysDataSource) Find(filter entity.DayFilter, joins ...entity.DayJoin) ([]*entity.Day, error) {
	days := []*entity.Day{}
	q, args := dayFilter(filter).generate()
	if err := ds.DB.Select(&days, q, args...); err != nil {
		return nil, catchErr("find days", err)
	}
	return days, database.ExecuteDaysPipeline(days, joins...)
}

func catchErr(msg string, err error) error {
	if err == nil {
		return nil
	}

	e := err.Error()
	logape.DefaultLogger.Debug(msg + ": " + e)
	switch e {
	case sql.ErrNoRows.Error():
		return fmt.Errorf("%w %s", database.ErrNotFound, msg)
	default:
		if strings.Contains(err.Error(), "UNIQUE constraint failed:") {
			return fmt.Errorf("%w %s", database.ErrUniqueCheckFailed, msg)
		}
		return errors.NewFatal(err.Error())
	}
}

func dayFilter(f entity.DayFilter) *sqlQueryBuilder {
	b := newBuilder("days")

	if len(f.IDs) > 0 {
		b.AddFilter("id", intToAny(f.IDs))
	}

	if len(f.Dates) > 0 {
		b.AddFilter("date", dateToAny(f.Dates))
	}

	return b
}

func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
