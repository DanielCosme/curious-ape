package sqlite

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/jmoiron/sqlx"
)

type HabitsDataSource struct {
	DB *sqlx.DB
}

func NewHabitsDataSource(db *sqlx.DB) *HabitsDataSource {
	return &HabitsDataSource{db}
}

func (ds HabitsDataSource) GetByUUID(id string) (*entity.Habit, error) {
	var habit *entity.Habit
	query := `SELECT * FROM habits where uuid=$1`
	return habit, ds.DB.Get(habit, query, id)
}

func (ds HabitsDataSource) Create(h *entity.Habit) error {
	query := `
		INSERT INTO habits ( id, state, time, creation_time, update_time ) 
		VALUES (:id, :state, :time, :creation_time, :update_time)
		RETURNING id
	`
	_, err := ds.DB.NamedExec(query, h)
	return err
}

func (ds HabitsDataSource) Find(filter *entity.HabitQuery) ([]*entity.Habit, error) {
	habits := []*entity.Habit{}
	query := `SELECT * from habits`
	return habits, ds.DB.Select(&habits, query)
}

func (ds HabitsDataSource) Update(habit *entity.Habit) (*entity.Habit, error) {
	panic("implement me")
}

func (ds HabitsDataSource) Delete(id string) error {
	panic("implement me")
}

func (ds HabitsDataSource) CreateHistoryEntry(hhe *entity.HabitHistoryEntry) error {
	panic("implement me")
}

func (ds HabitsDataSource) CreateCustomHabit(habitType *entity.HabitType) error {
	panic("implement me")
}
