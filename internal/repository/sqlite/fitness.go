package sqlite

import (
	"github.com/danielcosme/curious-ape/internal/core/database"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/jmoiron/sqlx"
)

type FitnessLogDataSource struct {
	DB *sqlx.DB
}

func (ds FitnessLogDataSource) Create(log *entity.FitnessLog) error {
	q := `
		INSERT INTO fitness_logs (
			day_id,
			date,
			type,
			start_time,
			end_time,
			title,
			origin,
			note,
			raw
	)	
		VALUES (
			:day_id,
			:date,
			:type,
			:start_time,
			:end_time,
			:title,
			:origin,
			:note,
			:raw
	)`
	res, err := ds.DB.NamedExec(q, log)
	if err != nil {
		return catchErr("create fitness log", err)
	}
	id, _ := res.LastInsertId()
	log.ID = int(id)
	return nil
}

func (ds FitnessLogDataSource) Update(log *entity.FitnessLog, join ...entity.FitnessLogJoin) (*entity.FitnessLog, error) {
	q := `
		UPDATE fitness_logs 
		SET 
			day_id = :day_id,
			date = :date,
			start_time = :start_time,
			end_time = :end_time,
			title = :title, 
			origin = :origin,
			note = :note, raw = :raw
		WHERE id = :id
	`
	_, err := ds.DB.NamedExec(q, log)
	if err != nil {
		return nil, catchErr("update fitness log", err)
	}
	return ds.Get(entity.FitnessLogFilter{ID: []int{log.ID}})
}

func (ds FitnessLogDataSource) Get(filter entity.FitnessLogFilter, joins ...entity.FitnessLogJoin) (*entity.FitnessLog, error) {
	fl := &entity.FitnessLog{}
	query, args := fitnessLogFilter(filter).generate()
	if err := ds.DB.Get(fl, query, args...); err != nil {
		return nil, catchErr("get fitness log", err)
	}
	return fl, catchErr("get fitness log", database.ExecuteFitnessLogPipeline([]*entity.FitnessLog{fl}, joins...))
}

func (ds FitnessLogDataSource) Find(filter entity.FitnessLogFilter, joins ...entity.FitnessLogJoin) ([]*entity.FitnessLog, error) {
	fls := []*entity.FitnessLog{}
	query, args := fitnessLogFilter(filter).generate()
	if err := ds.DB.Select(&fls, query, args...); err != nil {
		return nil, catchErr("find fitness log", err)
	}
	return fls, catchErr("find fitness log", database.ExecuteFitnessLogPipeline(fls, joins...))
}

func (ds FitnessLogDataSource) Delete(id int) error {
	_, err := ds.DB.Exec("DELETE FROM fitness_logs WHERE id = ?", id)
	return catchErr("delete fitness log", err)
}

func fitnessLogFilter(f entity.FitnessLogFilter) *sqlQueryBuilder {
	b := newBuilder("fitness_logs")

	if len(f.ID) > 0 {
		b.AddFilter("id", intToAny(f.ID))
	}

	if len(f.DayID) > 0 {
		b.AddFilter("day_id", intToAny(f.DayID))
	}

	if len(f.Date) > 0 {
		b.AddFilter("date", dateToAny(f.Date))
	}

	return b
}
