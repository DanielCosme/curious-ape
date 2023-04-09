package sqlite

import (
	"github.com/danielcosme/curious-ape/internal/core/database"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/jmoiron/sqlx"
)

type HabitsDataSource struct {
	DB *sqlx.DB
}

func (ds *HabitsDataSource) Create(h *entity.Habit) error {
	query := `
		INSERT INTO habits (
			day_id,
			habit_category_id,
			status
		) 
		VALUES (
			:day_id,
			:habit_category_id,
			:status
		)`
	res, err := ds.DB.NamedExec(query, h)
	if err != nil {
		return catchErr(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	h.ID = int(id)
	return nil
}
func (ds *HabitsDataSource) Update(data *entity.Habit, joins ...entity.HabitJoin) (*entity.Habit, error) {
	query := `
		UPDATE 	habits
		SET 	status = :status
		WHERE 	id = :id
	`
	_, err := ds.DB.NamedExec(query, data)
	if err != nil {
		return nil, catchErr(err)
	}
	return ds.Get(entity.HabitFilter{ID: []int{data.ID}}, joins...)
}

func (ds *HabitsDataSource) Get(filter entity.HabitFilter, joins ...entity.HabitJoin) (*entity.Habit, error) {
	habit := new(entity.Habit)
	query, args := habitFilter(filter).generate()
	if err := ds.DB.Get(habit, query, args...); err != nil {
		return nil, catchErr(err)
	}

	return habit, catchErr(database.ExecuteHabitsPipeline([]*entity.Habit{habit}, joins...))
}

func (ds *HabitsDataSource) Find(filter entity.HabitFilter, joins ...entity.HabitJoin) ([]*entity.Habit, error) {
	habits := []*entity.Habit{}
	query, args := habitFilter(filter).generate()
	if err := ds.DB.Select(&habits, query, args...); err != nil {
		return nil, catchErr(err)
	}
	return habits, catchErr(database.ExecuteHabitsPipeline(habits, joins...))
}

func (ds *HabitsDataSource) Delete(id int) error {
	_, err := ds.DB.Exec(`DELETE FROM habits WHERE id = ?`, id)
	return catchErr(err)
}

func (ds *HabitsDataSource) FindHabitCategories(filter entity.HabitCategoryFilter) ([]*entity.HabitCategory, error) {
	cs := []*entity.HabitCategory{}
	query, args := habitCategoryFilter(filter).generate()
	if err := ds.DB.Select(&cs, query, args...); err != nil {
		return nil, catchErr(err)
	}
	return cs, nil
}

func (ds *HabitsDataSource) GetHabitCategory(filter entity.HabitCategoryFilter) (*entity.HabitCategory, error) {
	hc := new(entity.HabitCategory)
	query, args := habitCategoryFilter(filter).generate()
	if err := ds.DB.Get(hc, query, args...); err != nil {
		return nil, catchErr(err)
	}
	return hc, nil
}

func (ds *HabitsDataSource) CreateHabitLog(hl *entity.HabitLog) error {
	query := `
		INSERT INTO habit_logs (
			habit_id,
			origin,
			is_automated,
			success,
			note
		) 
		VALUES (
			:habit_id,
			:origin,
			:is_automated,
			:success,
			:note
		)`
	res, err := ds.DB.NamedExec(query, hl)
	if err != nil {
		return catchErr(err)
	}
	id, _ := res.LastInsertId()
	hl.ID = int(id)
	return nil
}

func (ds *HabitsDataSource) UpdateHabitLog(data *entity.HabitLog) (*entity.HabitLog, error) {
	query := `
		UPDATE habit_logs
		SET 
			success = :success,
			origin = :origin,
			is_automated = :is_automated,
			note = :note
		WHERE id = :id
	`
	_, err := ds.DB.NamedExec(query, data)
	if err != nil {
		return nil, catchErr(err)
	}
	return ds.GetHabitLog(entity.HabitLogFilter{ID: []int{data.ID}})
}

func (ds *HabitsDataSource) GetHabitLog(filter entity.HabitLogFilter) (*entity.HabitLog, error) {
	hl := &entity.HabitLog{}
	query, args := habitLogFilter(filter).generate()
	if err := ds.DB.Get(hl, query, args...); err != nil {
		return nil, catchErr(err)
	}
	return hl, nil
}

func (ds *HabitsDataSource) FindHabitLogs(filter entity.HabitLogFilter) ([]*entity.HabitLog, error) {
	hls := []*entity.HabitLog{}
	query, args := habitLogFilter(filter).generate()
	if err := ds.DB.Select(&hls, query, args...); err != nil {
		return nil, catchErr(err)
	}
	return hls, nil
}

func (ds *HabitsDataSource) DeleteHabitLog(id int) error {
	_, err := ds.DB.Exec(`DELETE FROM habit_logs WHERE id = ?`, id)
	return catchErr(err)
}

func habitFilter(f entity.HabitFilter) *sqlQueryBuilder {
	b := newBuilder("habits")

	if len(f.ID) > 0 {
		b.AddFilter("id", intToInterface(f.ID))
	}

	if len(f.DayID) > 0 {
		b.AddFilter("day_id", intToInterface(f.DayID))
	}

	if len(f.CategoryID) > 0 {
		b.AddFilter("habit_category_id", intToInterface(f.CategoryID))
	}

	return b
}

func habitCategoryFilter(f entity.HabitCategoryFilter) *sqlQueryBuilder {
	b := newBuilder("habit_categories")

	if len(f.ID) > 0 {
		b.Data = append(b.Data, filterData{columnName: "id", values: intToInterface(f.ID)})
	}

	if len(f.Type) > 0 {
		b.Data = append(b.Data, filterData{columnName: "type", values: habitTypeToInterface(f.Type)})
	}

	return b
}

func habitLogFilter(f entity.HabitLogFilter) *sqlQueryBuilder {
	b := newBuilder("habit_logs")

	if len(f.ID) > 0 {
		b.AddFilter("id", intToInterface(f.ID))
	}

	if len(f.HabitID) > 0 {
		b.AddFilter("habit_id", intToInterface(f.HabitID))
	}

	if len(f.Origin) > 0 {
		values := make([]interface{}, len(f.Origin))
		for i, v := range f.Origin {
			values[i] = v
		}
		b.AddFilter("origin", values)
	}

	return b
}
