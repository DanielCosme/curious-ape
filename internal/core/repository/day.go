package repository

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
)

type Day interface {
	Create(*entity.Day) error
	Get(entity.DayFilter, ...entity.DayJoin) (*entity.Day, error)
	Find(entity.DayFilter, ...entity.DayJoin) ([]*entity.Day, error)
	// Helpers
	ToIDs([]*entity.Day) []int
}

func ExecuteDaysPipeline(days []*entity.Day, joins ...entity.DayJoin) error {
	if !(len(days) > 0) {
		return nil
	}

	for _, j := range joins {
		if err := j(days); err != nil {
			return err
		}
	}
	return nil
}
