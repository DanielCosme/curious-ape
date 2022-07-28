package sqlite

import (
	"github.com/danielcosme/curious-ape/internal/core/database"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/jmoiron/sqlx"
)

type SleepLogDataSource struct {
	DB *sqlx.DB
}

func (ds SleepLogDataSource) Create(log *entity.SleepLog) error {
	q := `
		INSERT INTO sleep_logs (day_id, date, start_time, end_time, is_main_sleep, is_automated, origin, total_time_in_bed,
		                        minutes_asleep, minutes_deep, minutes_rem, minutes_light, minutes_awake, raw)	
		values (:day_id, :date, :start_time, :end_time, :is_main_sleep, :is_automated, :origin, :total_time_in_bed,
		                        :minutes_asleep, :minutes_deep, :minutes_rem, :minutes_light, :minutes_awake, :raw) `
	res, err := ds.DB.NamedExec(q, log)
	if err != nil {
		return catchErr(err)
	}
	id, _ := res.LastInsertId()
	log.ID = int(id)
	return nil
}

func (ds SleepLogDataSource) Update(log *entity.SleepLog, joins ...entity.SleepLogJoin) (*entity.SleepLog, error) {
	q := `
		 UPDATE sleep_logs 
		 SET day_id = :day_id, date = :date, start_time = :start_time, end_time = :end_time, is_main_sleep = :is_main_sleep, 
		     is_automated = :is_automated, origin = :origin, total_time_in_bed = :total_time_in_bed, 
			 minutes_asleep = :minutes_asleep, minutes_rem = :minutes_rem, minutes_light = :minutes_light, 
		     minutes_awake = :minutes_awake, raw = :raw	
		WHERE id = :id
	`
	_, err := ds.DB.NamedExec(q, log)
	if err != nil {
		return nil, catchErr(err)
	}
	return ds.Get(entity.SleepLogFilter{ID: []int{log.ID}}, joins...)
}

func (ds SleepLogDataSource) Get(filter entity.SleepLogFilter, joins ...entity.SleepLogJoin) (*entity.SleepLog, error) {
	sl := &entity.SleepLog{}
	query, args := sleepLogFilter(filter).generate()
	if err := ds.DB.Get(sl, query, args...); err != nil {
		return nil, catchErr(err)
	}
	return sl, catchErr(database.ExecuteSleepLogPipeline([]*entity.SleepLog{sl}, joins...))
}

func (ds SleepLogDataSource) Find(filter entity.SleepLogFilter, joins ...entity.SleepLogJoin) ([]*entity.SleepLog, error) {
	sls := []*entity.SleepLog{}
	query, args := sleepLogFilter(filter).generate()
	if err := ds.DB.Select(&sls, query, args...); err != nil {
		return nil, catchErr(err)
	}
	return sls, catchErr(database.ExecuteSleepLogPipeline(sls, joins...))
}

func (ds SleepLogDataSource) Delete(id int) error {
	_, err := ds.DB.Exec("DELETE FROM sleep_logs WHERE id = ?", id)
	return catchErr(err)
}

func sleepLogFilter(f entity.SleepLogFilter) *sqlQueryBuilder {
	b := newBuilder("sleep_logs")

	if len(f.ID) > 0 {
		b.AddFilter("id", intToInterface(f.ID))
	}

	if len(f.DayID) > 0 {
		b.AddFilter("day_id", intToInterface(f.DayID))
	}

	return b
}
