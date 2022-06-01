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
	var args []interface{}
	query := `SELECT * FROM habits `

	if filter.DayID > 0 {
		query += fmt.Sprintf("WHERE day_id = ?")
		args = append(args, filter.DayID)
		if filter.CategoryID > 0 {
			query += fmt.Sprintf(" AND habit_category_id = ?")
			args = append(args, filter.CategoryID)
		}
	} else if filter.ID > 0 {
		query += fmt.Sprintf("WHERE id = ?")
		args = append(args, filter.ID)
	}

	if err := ds.DB.Get(habit, query, args...); err != nil {
		return nil, catchErr(err)
	}

	return habit, catchErr(repository.ExecuteHabitsPipeline([]*entity.Habit{habit}, joins...))
}

func (ds *HabitsDataSource) Create(h *entity.Habit, joins ...entity.HabitJoin) error {
	query := `
		INSERT INTO habits (day_id, habit_category_id, status) 
		VALUES (:day_id, :habit_category_id, :status)
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
		UPDATE 	habits
		SET 	status = :status
		WHERE 	id = :id
	`
	_, err := ds.DB.NamedExec(query, data)
	if err != nil {
		return nil, err
	}
	return ds.Get(entity.HabitFilter{ID: data.ID}, joins...)
}

func (ds *HabitsDataSource) Delete(id int) error {
	q := `DELETE FROM habits WHERE id = ?`
	_, err := ds.DB.Exec(q, id)
	return err
}

func (ds *HabitsDataSource) FindHabitCategories(filter entity.HabitFilter) ([]*entity.HabitCategory, error) {
	cs := []*entity.HabitCategory{}
	var args []interface{}
	query := `SELECT * FROM habit_categories`

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
	return hc, catchErr(ds.DB.Get(hc, "SELECT * FROM habit_categories WHERE id=?", filter.CategoryID))
}

func (ds *HabitsDataSource) CreateHabitLog(hl *entity.HabitLog) error {
	query := `
		INSERT INTO habit_logs (habit_id, origin, is_automated, success, note) 
		VALUES (:habit_id, :origin, :is_automated, :success, :note)
	`
	_, err := ds.DB.NamedExec(query, hl)
	return err
}

func (ds *HabitsDataSource) UpdateHabitLog(data *entity.HabitLog) (*entity.HabitLog, error) {
	query := `
		UPDATE habit_logs
		SET success = :success, origin = :origin, is_automated = :is_automated, note = :note
		WHERE id = :id
	`
	_, err := ds.DB.NamedExec(query, data)
	if err != nil {
		return nil, err
	}
	return ds.GetHabitLog(entity.HabitFilter{ID: data.ID})
}

func (ds *HabitsDataSource) GetHabitLog(filter entity.HabitFilter) (*entity.HabitLog, error) {
	hl := &entity.HabitLog{}
	query := `SELECT * FROM habit_logs `

	if filter.Origin != "" {
		query = fmt.Sprintf("%s WHERE id = ? AND origin = ?", query)
		return hl, catchErr(ds.DB.Get(hl, query, filter.ID, filter.Origin))
	}

	return hl, catchErr(ds.DB.Get(hl, query, filter.ID))
}

func (ds *HabitsDataSource) FindHabitLogs(filter entity.HabitFilter) ([]*entity.HabitLog, error) {
	hl := []*entity.HabitLog{}
	var args []interface{}
	query := `SELECT * FROM habit_logs`

	if len(filter.IDs) > 0 {
		var err error
		query, args, err = sqlx.In(fmt.Sprintf("%s WHERE habit_id IN (?)", query), filter.IDs)
		if err != nil {
			return nil, err
		}
	}

	return hl, ds.DB.Select(&hl, query, args...)
}

func (ds *HabitsDataSource) DeleteHabitLog(id int) error {
	_, err := ds.DB.Exec(`DELETE FROM habit_logs WHERE id = ?`, id)
	return err
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

func (ds *HabitsDataSource) ToIDs(hs []*entity.Habit) []int {
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
