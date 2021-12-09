package repository

import "github.com/danielcosme/curious-ape/internal/core/entity"

type Habit interface {
	GetByID(id entity.ID) (*entity.Habit, error)
	Create(habit *entity.Habit) error
	Find(query entity.HabitQuery) ([]*entity.Habit, error)
	Update(habit *entity.Habit) (*entity.Habit, error)
	Delete(id entity.ID) error
	CreateHistoryEntry(hhe *entity.HabitHistoryEntry) error
	CreateCustomHabit(habitType *entity.HabitType) error
}
