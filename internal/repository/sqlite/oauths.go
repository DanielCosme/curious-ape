package sqlite

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/jmoiron/sqlx"
)

type Oauth2DataSource struct {
	DB *sqlx.DB
}

func (ds *Oauth2DataSource) Create(o *entity.Oauth2) error {
	q := `
		INSERT INTO oauths (
			provider,
			access_token,
			refresh_token,
			expiration,
			type,
			toggl_workspace_id, 
			toggl_organization_id,
			toggl_project_ids
		)
		values (
			:provider,
			:access_token,
			:refresh_token,
			:expiration,
			:type,
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

func (ds *Oauth2DataSource) Update(o *entity.Oauth2) (*entity.Oauth2, error) {
	q := `
		UPDATE  oauths 
		SET 
			access_token = :access_token,
			refresh_token = :refresh_token,
			expiration = :expiration,
			type = :type, 
		    toggl_workspace_id = :toggl_workspace_id,
			toggl_organization_id = :toggl_organization_id,
			toggl_project_ids = :toggl_project_ids
		WHERE id = :id
	`
	_, err := ds.DB.NamedExec(q, o)
	if err != nil {
		return nil, catchErr(err)
	}
	return ds.Get(entity.Oauth2Filter{ID: []int{o.ID}})
}

func (ds *Oauth2DataSource) Get(filter entity.Oauth2Filter) (*entity.Oauth2, error) {
	o := new(entity.Oauth2)
	query, args := oauthFilter(filter).generate()
	if err := ds.DB.Get(o, query, args...); err != nil {
		return nil, catchErr(err)
	}
	return o, nil
}

func (ds *Oauth2DataSource) Find(filter entity.Oauth2Filter) ([]*entity.Oauth2, error) {
	oauths := []*entity.Oauth2{}
	query, args := oauthFilter(filter).generate()
	if err := ds.DB.Select(oauths, query, args...); err != nil {
		return nil, catchErr(err)
	}
	return oauths, nil
}

func (ds *Oauth2DataSource) Delete(id int) error {
	_, err := ds.DB.Exec("DELETE FROM oauths WHERE id = ?", id)
	return catchErr(err)
}

func oauthFilter(f entity.Oauth2Filter) *sqlQueryBuilder {
	b := newBuilder("oauths")

	if len(f.ID) > 0 {
		b.AddFilter("id", intToInterface(f.ID))
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
