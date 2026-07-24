package user

import (
	"context"

	"danicos.dev/daniel/curious-ape/pkg/gen/bob/models"
	"danicos.dev/daniel/curious-ape/pkg/persistence"
	"github.com/aarondl/opt/omit"
	"github.com/stephenafamo/bob"
)

type Params struct {
	ID       int64
	Username string
}

func Get(db bob.DB, p Params) (*User, error) {
	q := models.Users.Query()
	if p.ID > 0 {
		q.Apply(models.SelectWhere.Users.ID.EQ(int64(p.ID)))
	}
	if p.Username != "" {
		q.Apply(models.SelectWhere.Users.Username.EQ(p.Username))
	}
	m, err := q.One(context.Background(), db)
	if err != nil {
		return nil, persistence.CatchDBErr("app: user get:", err)
	}
	u := &User{Name: m.Username, Password: m.Password}
	u.ID = m.ID
	return u, nil
}

func create(db bob.DB, username, hashedPassword string) error {
	setter := &models.UserSetter{
		Username: omit.From(username),
		Password: omit.From(string(hashedPassword)),
		// TODO: remove this from the database, via migration.
		Role:  omit.From("admin"),
		Email: omit.From("danicosme@pm.me"),
	}
	_, err := models.Users.Insert(setter).One(context.Background(), db)
	return persistence.CatchDBErr("app: user create", err)
}
