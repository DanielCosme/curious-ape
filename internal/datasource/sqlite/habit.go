package sqlite

import (
	"fmt"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/core/repository"
	"github.com/jmoiron/sqlx"
)

type HabitsDataSource struct {
	DB *sqlx.DB
}

func (ds *HabitsDataSource) Get(filter entity.HabitFilter, joins ...entity.HabitJoin) (*entity.Habit, error) {
	habit := new(entity.Habit)
	query, args := newHabitQueryBuilder(filter).Generate()
	if err := ds.DB.Get(habit, query, args...); err != nil {
		return nil, err
	}
	return habit, repository.ExecuteHabitsPipeline([]*entity.Habit{habit}, joins...)
}

func (ds *HabitsDataSource) Create(h *entity.Habit, joins ...entity.HabitJoin) error {
	query := `
		INSERT INTO habits (day_id, habit_category_id, success, origin, is_automated, note) 
		VALUES (:day_id, :habit_category_id, :success, :origin, :is_automated, :note)
	`
	res, err := ds.DB.NamedExec(query, h)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	h.ID = int(id)
	return repository.ExecuteHabitsPipeline([]*entity.Habit{h}, joins...)
}

func (ds *HabitsDataSource) Find(filter entity.HabitFilter, joins ...entity.HabitJoin) ([]*entity.Habit, error) {
	habits := []*entity.Habit{}
	query := `SELECT * from habits`

	if len(filter.DayIDs) > 0 {
		q, args, err := sqlx.In(fmt.Sprintf("%s WHERE day_id IN (?)", query), filter.DayIDs)
		if err != nil {
			return nil, err
		}

		if err := ds.DB.Select(&habits, q, args...); err != nil {
			return nil, err
		}
	} else {
		if err := ds.DB.Select(&habits, query); err != nil {
			return nil, err
		}
	}

	return habits, repository.ExecuteHabitsPipeline(habits, joins...)
}

func (ds *HabitsDataSource) Update(data *entity.Habit, joins ...entity.HabitJoin) (*entity.Habit, error) {
	query := `
		UPDATE habits
		SET success = :success, origin = :origin, is_automated = :is_automated, note = :note
		WHERE id = :id
	`
	_, err := ds.DB.NamedExec(query, data)
	if err != nil {
		return nil, err
	}
	return ds.Get(entity.HabitFilter{ID: []int{data.ID}}, joins...)
}

func (ds *HabitsDataSource) Delete(id int) error {
	q := `DELETE FROM habits WHERE id = ?`
	_, err := ds.DB.Exec(q, id)
	return err
}

func (ds *HabitsDataSource) FindHabitCategories(filter entity.HabitFilter) ([]*entity.HabitCategory, error) {
	cs := []*entity.HabitCategory{}
	query, args := newHabitCategoryQueryBuilder().Generate()

	if len(filter.CategoryIDs) > 0 {
		q, args, err := sqlx.In(fmt.Sprintf("%s WHERE id IN (?)", query), filter.CategoryIDs)
		if err != nil {
			return nil, err
		}

		return cs, ds.DB.Select(&cs, q, args...)
	}

	return cs, ds.DB.Select(&cs, query, args...)
}

func (ds *HabitsDataSource) GetHabitCategory(filter entity.HabitFilter) (*entity.HabitCategory, error) {
	hc := new(entity.HabitCategory)
	// query, args := newHabitCategoryQueryBuilder().Generate()
	return hc, parseError(ds.DB.Get(hc, "SELECT * FROM habit_categories WHERE id=?", filter.ID[0]))
}

func newHabitCategoryQueryBuilder() *QueryBuilder {
	q := &QueryBuilder{tableName: "habit_categories"}
	return q
}

func newHabitQueryBuilder(f entity.HabitFilter) *QueryBuilder {
	q := &QueryBuilder{tableName: "habits"}

	for _, v := range f.ID {
		q.Add("id", v)
	}

	if len(q.params) > 0 {
		q.where = true
	}
	return q
}

func (ds *HabitsDataSource) ToDayIDs(hs []*entity.Habit) []int {
	dayIDs := []int{}
	dayIDsMap := map[int]int{}
	for _, h := range hs {
		if _, ok := dayIDsMap[h.DayID]; !ok {
			dayIDs = append(dayIDs, h.DayID)
			dayIDsMap[h.DayID] = h.DayID
		}
	}
	return dayIDs
}

func (ds *HabitsDataSource) ToCategoryIDs(hs []*entity.Habit) []int {
	categoryIDs := []int{}
	categoryIDsMap := map[int]int{}
	for _, h := range hs {
		if _, ok := categoryIDsMap[h.CategoryID]; !ok {
			categoryIDs = append(categoryIDs, h.CategoryID)
			categoryIDsMap[h.CategoryID] = h.CategoryID
		}
	}
	return categoryIDs
}
