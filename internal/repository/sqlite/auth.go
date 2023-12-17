package sqlite

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/jmoiron/sqlx"
)

type AuthenticationDataSource struct {
	DB *sqlx.DB
}

func (ds *AuthenticationDataSource) Create(o *entity.Auth) error {
	q := `
		INSERT INTO auths (
			provider,
			access_token,
			refresh_token,
			expiration,
			token_type,
			toggl_workspace_id, 
			toggl_organization_id,
			toggl_project_ids
		)
		values (
			:provider,
			:access_token,
			:refresh_token,
			:expiration,
			:token_type,
			:toggl_workspace_id,
			:toggl_organization_id,
			:toggl_project_ids
		)`
	res, err := ds.DB.NamedExec(q, o)
	if err != nil {
		return catchErr(err)
	}
	id, _ := res.LastInsertId()
	o.ID = int(id)
	return nil
}

func (ds *AuthenticationDataSource) Update(o *entity.Auth) (*entity.Auth, error) {
	q := `
		UPDATE  auths 
		SET 
			access_token = :access_token,
			refresh_token = :refresh_token,
			expiration = :expiration,
			token_type = :token_type, 
		    toggl_workspace_id = :toggl_workspace_id,
			toggl_organization_id = :toggl_organization_id,
			toggl_project_ids = :toggl_project_ids
		WHERE id = :id
	`
	_, err := ds.DB.NamedExec(q, o)
	if err != nil {
		return nil, catchErr(err)
	}
	return ds.Get(entity.AuthFilter{ID: []int{o.ID}})
}

func (ds *AuthenticationDataSource) Get(filter entity.AuthFilter) (*entity.Auth, error) {
	o := new(entity.Auth)
	query, args := authFilter(filter).generate()
	if err := ds.DB.Get(o, query, args...); err != nil {
		return nil, catchErr(err)
	}
	return o, nil
}

func (ds *AuthenticationDataSource) Find(filter entity.AuthFilter) ([]*entity.Auth, error) {
	auths := []*entity.Auth{}
	query, args := authFilter(filter).generate()
	if err := ds.DB.Select(auths, query, args...); err != nil {
		return nil, catchErr(err)
	}
	return auths, nil
}

func (ds *AuthenticationDataSource) Delete(id int) error {
	_, err := ds.DB.Exec("DELETE FROM auths WHERE id = ?", id)
	return catchErr(err)
}

func authFilter(f entity.AuthFilter) *sqlQueryBuilder {
	b := newBuilder("auths")

	if len(f.ID) > 0 {
		b.AddFilter("id", intToAny(f.ID))
	}

	if len(f.Provider) > 0 {
		values := make([]interface{}, len(f.Provider))
		for i, v := range f.Provider {
			values[i] = v
		}
		b.AddFilter("provider", values)
	}

	return b
}
