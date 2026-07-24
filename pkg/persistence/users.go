package persistence

import (
	"context"

	"danicos.dev/daniel/curious-ape/pkg/gen/bob/models"
	"github.com/stephenafamo/bob"
)

type Users struct {
	db bob.DB
}

func (u *Users) Create(s *models.UserSetter) (*models.User, error) {
	return models.Users.Insert(s).One(context.Background(), u.db)
}

func (u *Users) Exists(id int) (bool, error) {
	return models.UserExists(context.Background(), u.db, int64(id))
}

type UserParams struct {
	ID       int
	Username string
}

func (u *Users) Get(f UserParams) (*models.User, error) {
	q := models.Users.Query()
	if f.ID > 0 {
		q.Apply(models.SelectWhere.Users.ID.EQ(int64(f.ID)))
	}
	if f.Username != "" {
		q.Apply(models.SelectWhere.Users.Username.EQ(f.Username))
	}
	m, err := q.One(context.Background(), u.db)
	return m, CatchDBErr("GET USER", err)
}
