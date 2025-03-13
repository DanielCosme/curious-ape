package database

import (
	"context"
	"github.com/aarondl/opt/omit"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/database/gen/models"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
)

type Auths struct {
	db bob.DB
}

func (a *Auths) Upsert(s *models.AuthSetter) (*models.Auth, error) {
	auth, err := models.Auths.Insert(s).One(context.Background(), a.db)
	if err == nil {
		return auth, nil

	}

	if models.AuthErrors.ErrUniqueProvider.Is(err) {
		auth, err = a.Get(AuthParams{Provider: core.Integration(s.Provider.GetOrZero())})
		if err != nil {
			return nil, err
		}

		s.ID = omit.From(auth.ID)
		return models.Auths.Update(s.UpdateMod()).One(context.Background(), a.db)
	}

	return nil, catchDBErr("auths: upsert", err)
}

func (a *Auths) Get(p AuthParams) (*models.Auth, error) {
	res, err := p.BuildQuery().One(context.Background(), a.db)
	return res, catchDBErr("get auth", err)
}

type AuthParams struct {
	Provider core.Integration
}

func (f AuthParams) BuildQuery() *sqlite.ViewQuery[*models.Auth, models.AuthSlice] {
	q := models.Auths.Query()
	if f.Provider != "" {
		q.Apply(models.SelectWhere.Auths.Provider.EQ(string(f.Provider)))
	}
	return q
}
