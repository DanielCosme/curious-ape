package database

import (
	"context"
	"github.com/danielcosme/curious-ape/internal/database/gen/models"
	"github.com/stephenafamo/bob"
)

type Auths struct {
	db bob.DB
}

func (a *Auths) Upsert(s *models.AuthSetter) (*models.Auth, error) {
	return models.Auths.Upsert(
		context.Background(),
		a.db,
		true,
		[]string{"provider"},
		nil,
		s,
	)
}

func (a *Auths) Get(p AuthParams) (*models.Auth, error) {
	res, err := p.BuildQuery(a.db).One()
	return res, catchErr("get auth", err)
}
