package persistence

import (
	"context"
	"git.danicos.dev/daniel/curious-ape/database/gen/dberrors"
	"git.danicos.dev/daniel/curious-ape/database/gen/models"
	"git.danicos.dev/daniel/curious-ape/pkg/core"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
)

type Auths struct {
	db bob.DB
}

func (a *Auths) Upsert(s *models.OauthTokenSetter) (*models.OauthToken, error) {
	auth, err := models.OauthTokens.Insert(s).One(context.Background(), a.db)
	if err == nil {
		return auth, nil
	}
	if dberrors.OauthTokenErrors.ErrUniqueSqliteAutoindexOauthToken1.Is(err) {
		return models.OauthTokens.
			Update(s.UpdateMod(), models.UpdateWhere.OauthTokens.Provider.EQ(s.Provider.GetOrZero())).
			One(context.Background(), a.db)
	}
	return nil, catchDBErr("auths: upsert", err)
}

func (a *Auths) Get(p AuthParams) (*models.OauthToken, error) {
	res, err := p.BuildQuery().One(context.Background(), a.db)
	return res, catchDBErr("get auth", err)
}

type AuthParams struct {
	Provider core.Integration
}

func (f AuthParams) BuildQuery() *sqlite.ViewQuery[*models.OauthToken, models.OauthTokenSlice] {
	q := models.OauthTokens.Query()
	if f.Provider != "" {
		q.Apply(models.SelectWhere.OauthTokens.Provider.EQ(string(f.Provider)))
	}
	return q
}
