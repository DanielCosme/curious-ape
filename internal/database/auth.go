package database

import (
	"github.com/danielcosme/curious-ape/internal/entity"
)

type Authentication interface {
	Create(*entity.Auth) error
	Update(*entity.Auth) (*entity.Auth, error)
	Get(entity.AuthFilter) (*entity.Auth, error)
	Find(entity.AuthFilter) ([]*entity.Auth, error)
	Delete(int) error
}
