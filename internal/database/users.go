package database

import (
	"github.com/danielcosme/curious-ape/internal/entity"
)

type User interface {
	Create(*entity.User) error
	Update(*entity.User) (*entity.User, error)
	Get(entity.UserFilter) (*entity.User, error)
	Delete(int) error
}
