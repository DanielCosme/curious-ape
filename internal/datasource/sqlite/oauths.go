package sqlite

import (
	"fmt"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/core/repository"
	"github.com/jmoiron/sqlx"
)

type Oauth2DataSource struct {
	DB *sqlx.DB
}

func (ds *Oauth2DataSource) Create(o *entity.Oauth2) error {
	q := `
		INSERT INTO oauths (provider, access_token, refresh_token, expiration, type)	
		values (:provider, :access_token, :refresh_token, :expiration, :type) `
	_, err := ds.DB.NamedExec(q, o)
	return err
}

func (ds *Oauth2DataSource) Update(o *entity.Oauth2) (*entity.Oauth2, error) {
	q := `
		UPDATE  oauths 
		SET access_token = :access_token, refresh_token = :refresh_token, expiration = :expiration, type = :type
		WHERE id = :id
	`
	_, err := ds.DB.NamedExec(q, o)
	if err != nil {
		return nil, err
	}
	return ds.Get(entity.Oauth2Filter{IDs: []int{o.ID}})
}

func (ds *Oauth2DataSource) Get(filter entity.Oauth2Filter) (*entity.Oauth2, error) {
	o := new(entity.Oauth2)
	q := `SELECT * FROM oauths`

	if len(filter.IDs) > 0 {
		q = fmt.Sprintf("%s %s",q, "WHERE id = ?")
		return o, parseError(ds.DB.Get(o, q, filter.IDs[0]))
	} else if len(filter.Providers) > 0 {
		q = fmt.Sprintf("%s %s",q, "WHERE provider = ?")
		return o, parseError(ds.DB.Get(o, q, filter.Providers[0]))
	}
	return nil, repository.ErrNotFound
}

func (ds *Oauth2DataSource) Find(filter entity.Oauth2Filter) ([]*entity.Oauth2, error) {
	o := []*entity.Oauth2{}
	q := `SELECT * FROM oauths`

	return o, ds.DB.Select(o, q)
}

func (ds *Oauth2DataSource) Delete(id int) error {
	_, err := ds.DB.Exec("DELETE FROM oauths WHERE id = ?", id)
	return err
}
