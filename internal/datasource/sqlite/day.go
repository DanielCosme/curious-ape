package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/core/repository"
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

func (ds *DaysDataSource) Get(filter entity.DayFilter) (*entity.Day, error) {
	day := new(entity.Day)
	q, args := newDayQueryBuilder(filter).Generate()
	return day, parseError(ds.DB.Get(day, q, args...))
}

func parseError(err error) error {
	if err == nil {
		return nil
	}
	switch err.Error() {
	case sql.ErrNoRows.Error():
		return repository.ErrNotFound
	default:
		return fmt.Errorf("%s %w", repository.ErrNotImplemented, err)
	}
}

func (ds *DaysDataSource) Find(filter entity.DayFilter) ([]*entity.Day, error) {
	days := []*entity.Day{}
	query := `SELECT * from "days"`
	if len(filter.IDs) > 0 {
		q, args, err := sqlx.In(fmt.Sprintf("%s WHERE id IN (?)", query), filter.IDs)
		if err != nil {
			return nil, err
		}

		return days, ds.DB.Select(&days, q, args...)
	}

	return days, ds.DB.Select(&days, query)
}

func newDayQueryBuilder(f entity.DayFilter) *QueryBuilder {
	q := &QueryBuilder{
		params:    []queryTouple{},
		tableName: "days",
	}

	for _, v := range f.IDs {
		q.Add("id", v)
	}
	for _, v := range f.Date {
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
