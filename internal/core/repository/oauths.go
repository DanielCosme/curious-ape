package repository

import "github.com/danielcosme/curious-ape/internal/core/entity"

type Oauth2 interface {
	Create(oauth *entity.Oauth2) error
	Update(oauth *entity.Oauth2) (*entity.Oauth2, error)
	Get(filter entity.Oauth2Filter) (*entity.Oauth2, error)
	Find(filter entity.Oauth2Filter) ([]*entity.Oauth2, error)
	Delete(id int) error
}
