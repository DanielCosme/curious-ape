package application

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/core/repository"
)

type HabitsInteractor struct {
	Service repository.Habit
}

func (hc *HabitsInteractor) Create(h *entity.Habit) (*entity.Habit, error) {
	// At least one history entry?
	if err := hc.Service.Create(h); err != nil {
		return nil, err
	}
	return h, nil
}

func (hc *HabitsInteractor) GetAll() ([]*entity.Habit, error) {
	q := entity.HabitQuery{
		Query:     nil,
		DateQuery: nil,
	}
	return hc.Service.Find(q)
}
