package repository

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
)

type Day interface {
	Create(day *entity.Day) error
	Get(filter entity.DayFilter) (*entity.Day, error)
	Find(filter entity.DayFilter) ([]*entity.Day, error)
}
