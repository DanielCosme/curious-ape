package application

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/core/repository"
)

type HabitsInteractor struct {
	repo repository.Habit
}

func (hi *HabitsInteractor) Create(h *entity.Habit) (*entity.Habit, error) {
	// At least one history entry?
	if err := hi.repo.Create(h); err != nil {
		return nil, err
	}
	return h, nil
}

func (hi *HabitsInteractor) GetAll() ([]*entity.Habit, error) {
	q := &entity.HabitQuery{
		Query:     nil,
		DateQuery: nil,
	}
	return hi.repo.Find(q)
}
