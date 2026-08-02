package integration

import (
	"context"
	"fmt"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/gen/bob/dberrors"
	"danicos.dev/daniel/curious-ape/pkg/gen/bob/models"
	"danicos.dev/daniel/curious-ape/pkg/persistence"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
)

func upsert(db bob.DB, s *models.OauthTokenSetter) (*models.OauthToken, error) {
	auth, err := models.OauthTokens.Insert(s).One(context.Background(), db)
	if err == nil {
		return auth, nil
	}
	if dberrors.OauthTokenErrors.ErrUniqueSqliteAutoindexOauthToken1.Is(err) {
		return models.OauthTokens.
			Update(s.UpdateMod(), models.UpdateWhere.OauthTokens.Provider.EQ(s.Provider.GetOrZero())).
			One(context.Background(), db)
	}
	return nil, persistence.CatchDBErr("auths: upsert", err)
}

func get(db bob.DB, p AuthParams) (*models.OauthToken, error) {
	res, err := p.buildQuery().One(context.Background(), db)
	return res, persistence.CatchDBErr(fmt.Sprintf("repository: get auth: %s", p.Provider), err)
}

type AuthParams struct {
	Provider core.Integration
}

func (f AuthParams) buildQuery() *sqlite.ViewQuery[*models.OauthToken, models.OauthTokenSlice] {
	q := models.OauthTokens.Query()
	if f.Provider != "" {
		q.Apply(models.SelectWhere.OauthTokens.Provider.EQ(string(f.Provider)))
	}
	return q
}
