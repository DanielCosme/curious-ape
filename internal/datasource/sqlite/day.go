package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/core/repository"
	"github.com/danielcosme/curious-ape/sdk/errors"
	"github.com/jmoiron/sqlx"
)

type DaysDataSource struct {
	DB *sqlx.DB
}

func (ds *DaysDataSource) Create(d *entity.Day) error {
	query := `
		INSERT INTO "days" ("date") 
		VALUES (:date);
	`
	_, err := ds.DB.NamedExec(query, d)
	return err
}

func (ds *DaysDataSource) Get(filter entity.DayFilter, joins ...entity.DayJoin) (*entity.Day, error) {
	day := new(entity.Day)
	var args []interface{}
	q := `SELECT * FROM days`

	if !filter.Date.IsZero() {
		q = fmt.Sprintf("%s WHERE date = ?", q)
		args = append(args, filter.Date)
	}

	return day, catchErr(ds.DB.Get(day, q, args...))
}

func (ds *DaysDataSource) Find(filter entity.DayFilter, joins ...entity.DayJoin) ([]*entity.Day, error) {
	days := []*entity.Day{}
	query := `SELECT * from "days"`
	if len(filter.IDs) > 0 {
		q, args, err := sqlx.In(fmt.Sprintf("%s WHERE id IN (?)", query), filter.IDs)
		if err != nil {
			return nil, err
		}

		if err := ds.DB.Select(&days, q, args...); err != nil {
			return nil, err
		}
	} else {
		if err := ds.DB.Select(&days, query); err != nil {
			return nil, err
		}
	}

	return days, repository.ExecuteDaysPipeline(days, joins...)
}

func (ds *DaysDataSource) ToIDs(days []*entity.Day) []int {
	ids := []int{}
	for _, d := range days {
		ids = append(ids, d.ID)
	}
	return ids
}

func catchErr(err error) error {
	if err == nil {
		return nil
	}
	switch err.Error() {
	case sql.ErrNoRows.Error():
		return repository.ErrNotFound
	default:
		return errors.NewFatal(err.Error())
	}
}
