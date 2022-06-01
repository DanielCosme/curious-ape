package repository

import "github.com/danielcosme/curious-ape/internal/core/entity"

type Habit interface {
	// habit
	Create(*entity.Habit, ...entity.HabitJoin) error
	Update(*entity.Habit, ...entity.HabitJoin) (*entity.Habit, error)
	Get(entity.HabitFilter, ...entity.HabitJoin) (*entity.Habit, error)
	Find(entity.HabitFilter, ...entity.HabitJoin) ([]*entity.Habit, error)
	Delete(id int) error
	// habit log
	CreateHabitLog(*entity.HabitLog) error
	UpdateHabitLog(*entity.HabitLog) (*entity.HabitLog, error)
	GetHabitLog(entity.HabitFilter) (*entity.HabitLog, error)
	FindHabitLogs(entity.HabitFilter) ([]*entity.HabitLog, error)
	DeleteHabitLog(id int) error
	// habit category
	GetHabitCategory(filter entity.HabitFilter) (*entity.HabitCategory, error)
	FindHabitCategories(filter entity.HabitFilter) ([]*entity.HabitCategory, error)
	// Helpers
	ToIDs(hs []*entity.Habit) []int
	ToDayIDs(hs []*entity.Habit) []int
	ToCategoryIDs(hs []*entity.Habit) []int
}

func ExecuteHabitsPipeline(hs []*entity.Habit, hjs ...entity.HabitJoin) error {
	if !(len(hs) > 0) {
		return nil
	}

	for _, hj := range hjs {
		if err := hj(hs); err != nil {
			return err
		}
	}
	return nil
}
