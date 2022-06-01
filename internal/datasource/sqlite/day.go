package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/core/repository"
	"github.com/danielcosme/curious-ape/sdk/errors"
	"github.com/jmoiron/sqlx"
	"strings"
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
	q, args := newDayQueryBuilder(filter).Generate()
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

func newDayQueryBuilder(f entity.DayFilter) *QueryBuilder {
	q := &QueryBuilder{
		params:    []queryTouple{},
		tableName: "days",
	}

	for _, v := range f.IDs {
		q.Add("id", v)
	}
	for _, v := range f.Dates {
		q.Add("date", v)
	}

	if len(q.params) > 0 {
		q.where = true
	}
	return q
}

type QueryBuilder struct {
	params    []queryTouple
	tableName string
	where     bool
}

type queryTouple struct {
	key   string
	value interface{}
}

func (qb *QueryBuilder) Add(key string, value interface{}) {
	qb.params = append(qb.params, queryTouple{key: key, value: value})
}

func (g *QueryBuilder) Generate() (string, []interface{}) {
	var args []interface{}
	query := fmt.Sprintf("SELECT * FROM %s", g.tableName)

	if g.where {
		query += " WHERE "
	}

	lines := []string{}
	for _, v := range g.params {
		line := fmt.Sprintf("%s = ?", v.key)
		lines = append(lines, line)
		args = append(args, v.value)
	}
	query += strings.Join(lines, " OR ")

	return query, args
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
		return entity.ErrNotFound
	default:
		return errors.NewFatal(err.Error())
	}
}
