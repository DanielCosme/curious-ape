package sqlite

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/core/repository"
	"github.com/jmoiron/sqlx"
)

type SleepLogDataSource struct {
	DB *sqlx.DB
}

func (ds SleepLogDataSource) Create(log *entity.SleepLog) error {
	q := `
		INSERT INTO sleep_logs (day_id, start_time, end_time, is_main_sleep, is_automated, origin, total_time_in_bed,
		                        minutes_asleep, minutes_deep, minutes_rem, minutes_light, minutes_awake, raw)	
		values (:day_id, :start_time, :end_time, :is_main_sleep, :is_automated, :origin, :total_time_in_bed,
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
		 SET day_id = :day_id, start_time = :start_time, end_time = :end_time, is_main_sleep = :is_main_sleep, 
		     is_automated = :is_automated, origin = :origin, total_time_in_bed = :total_time_in_bed, 
			 minutes_asleep = :minutes_asleep, minutes_rem = :minutes_rem, minutes_light = :minutes_light, 
		     minutes_awake = :minutes_awake, raw = :raw	
		WHERE id = :id
	`
	res, err := ds.DB.NamedExec(q, log)
	if err != nil {
		return nil, catchErr(err)
	}
	id, _ := res.LastInsertId()
	return ds.Get(entity.SleepLogFilter{ID: []int{int(id)}}, joins...)
}

func (ds SleepLogDataSource) Get(filter entity.SleepLogFilter, joins ...entity.SleepLogJoin) (*entity.SleepLog, error) {
	sl := &entity.SleepLog{}
	query, args := sleepLogFilter(filter).generate()
	if err := ds.DB.Get(sl, query, args...); err != nil {
		return nil, catchErr(err)
	}
	return sl, catchErr(repository.ExecuteSleepLogPipeline([]*entity.SleepLog{sl}, joins...))
}

func (ds SleepLogDataSource) Find(filter entity.SleepLogFilter, joins ...entity.SleepLogJoin) ([]*entity.SleepLog, error) {
	sls := []*entity.SleepLog{}
	query, args := sleepLogFilter(filter).generate()
	if err := ds.DB.Select(sls, query, args...); err != nil {
		return nil, catchErr(err)
	}
	return sls, catchErr(repository.ExecuteSleepLogPipeline(sls, joins...))
}

func (ds SleepLogDataSource) Delete(id int) error {
	_, err := ds.DB.Exec("DELETE FROM sleep_logs WHERE id = ?", id)
	return catchErr(err)
}

func sleepLogFilter(f entity.SleepLogFilter) *sqlBuilder {
	b := newBuilder("sleep_logs")

	if len(f.ID) > 0 {
		b.AddFilter("id", intToInterface(f.ID))
	}

	if len(f.DayID) > 0 {
		b.AddFilter("day_id", intToInterface(f.ID))
	}

	return b
}

func (ds *SleepLogDataSource) ToDayIDs(sls []*entity.SleepLog) []int {
	dayIDs := []int{}
	dayIDsMap := map[int]int{}
	for _, h := range sls {
		if _, ok := dayIDsMap[h.DayID]; !ok {
			dayIDs = append(dayIDs, h.DayID)
			dayIDsMap[h.DayID] = h.DayID
		}
	}
	return dayIDs
}

func (ds *SleepLogDataSource) ToIDs(hs []*entity.SleepLog) []int {
	IDs := []int{}
	mapHabitIDs := map[int]int{}
	for _, h := range hs {
		if _, ok := mapHabitIDs[h.ID]; !ok {
			IDs = append(IDs, h.ID)
			mapHabitIDs[h.ID] = h.ID
		}
	}
	return IDs
}
