package sqlite

import (
	"github.com/danielcosme/curious-ape/internal/core/database"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/jmoiron/sqlx"
)

type UsersDataSource struct {
	DB *sqlx.DB
}

var _ database.User = (*UsersDataSource)(nil)

func (ds *UsersDataSource) Create(u *entity.User) error {
	q := `
		INSERT INTO users (
			name,
			password,
			role,
			email
		)
		VALUES (
			:name,
			:password,
			:role,
			:email
		)`
	res, err := ds.DB.NamedExec(q, u)
	u.ID = lastInsertID(res)
	return err
}

func (ds *UsersDataSource) Update(u *entity.User) (*entity.User, error) {
	q := `
		UPDATE  users 
		SET 
		name = :name,
		password = :password,
		role = :role,
		email = :email
		WHERE id = :id
	`
	_, err := ds.DB.NamedExec(q, u)
	if err != nil {
		return nil, catchErr(err)
	}
	return ds.Get(entity.UserFilter{ID: u.ID})
}

func (ds *UsersDataSource) Get(filter entity.UserFilter) (*entity.User, error) {
	u := new(entity.User)
	query, args := userFilter(filter).generate()
	if err := ds.DB.Get(u, query, args...); err != nil {
		return nil, catchErr(err)
	}
	return u, nil
}

func (ds *UsersDataSource) Delete(id int) error {
	_, err := ds.DB.Exec("DELETE FROM users WHERE id = ?", id)
	return catchErr(err)
}

func userFilter(f entity.UserFilter) *sqlQueryBuilder {
	b := newBuilder("users")

	if f.ID > 0 {
		b.AddFilter("id", []any{f.ID})
	}

	if f.Name != "" {
		b.AddFilter("name", []any{f.Name})
	}

	if f.Role != "" {
		b.AddFilter("role", []any{f.Role})
	}

	if f.Password != "" {
		b.AddFilter("password", []any{f.Password})
	}

	return b
}
