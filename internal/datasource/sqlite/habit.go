package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/danielcosme/curious-ape/internal/core/entity"
)

type HabitService struct {
	db *sql.DB
}

func NewHabitsService(db *sql.DB) *HabitService {
	return &HabitService{db}
}

func (h HabitService) GetByID(id entity.UUID) (*entity.Habit, error) {
	panic("implement me")
}

func (h HabitService) Create(habit *entity.Habit) error {
	fmt.Println("We have created")
	return nil
}

func (h HabitService) Find(query *entity.HabitQuery) ([]*entity.Habit, error) {
	panic("implement me")
}

func (h HabitService) Update(habit *entity.Habit) (*entity.Habit, error) {
	panic("implement me")
}

func (h HabitService) Delete(id entity.UUID) error {
	panic("implement me")
}

func (h HabitService) CreateHistoryEntry(hhe *entity.HabitHistoryEntry) error {
	panic("implement me")
}

func (h HabitService) CreateCustomHabit(habitType *entity.HabitType) error {
	panic("implement me")
}
