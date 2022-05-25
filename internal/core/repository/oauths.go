package repository

import "github.com/danielcosme/curious-ape/internal/core/entity"

type Oauth2 interface {
	Create(*entity.Oauth2) error
	Update(*entity.Oauth2) (*entity.Oauth2, error)
	Get(entity.Oauth2Filter) (*entity.Oauth2, error)
	Find(entity.Oauth2Filter) ([]*entity.Oauth2, error)
	Delete(int) error
}
