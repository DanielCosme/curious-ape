package repository

import "github.com/danielcosme/curious-ape/internal/core/entity"

type Habit interface {
	Create(habit *entity.Habit, joins ...entity.HabitJoin) error
	Update(habit *entity.Habit, joins ...entity.HabitJoin) (*entity.Habit, error)
	Get(filter entity.HabitFilter, joins ...entity.HabitJoin) (*entity.Habit, error)
	Find(filter entity.HabitFilter, joins ...entity.HabitJoin) ([]*entity.Habit, error)
	Delete(id int) error
	GetHabitCategory(filter entity.HabitFilter) (*entity.HabitCategory, error)
	FindHabitCategories(filter entity.HabitFilter) ([]*entity.HabitCategory, error)
	// Helpers
	ToDayIDs(hs []*entity.Habit) []int
	ToCategoryIDs(hs []*entity.Habit) []int
}

func ExecuteHabitsPipeline(hs []*entity.Habit, hjs ...entity.HabitJoin) error {
	for _, hj := range hjs {
		if err := hj(hs); err != nil {
			return err
		}
	}
	return nil
}
