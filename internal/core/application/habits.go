package application

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/core/repository"
	"time"
)

type HabitsInteractor struct {
	db repository.Habit
}

func (hi *HabitsInteractor) Create(h *entity.Habit) (*entity.Habit, error) {
	// At least one history entry?
	t := time.Now()
	h.Entity = entity.NewEntity()
	h.Time = t

	if err := hi.db.Create(nil); err != nil {
		return nil, err
	}
	return h, nil
}

func (hi *HabitsInteractor) GetAll() ([]*entity.Habit, error) {
	return hi.db.Find(nil)
}

func (hi *HabitsInteractor) GetByID(id string) (*entity.Habit, error) {
	return hi.db.GetByUUID(id)
}
