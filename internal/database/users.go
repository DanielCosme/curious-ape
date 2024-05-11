package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database/gen/models"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"log/slog"
)

type UserF struct {
	Role     core.Role
	Username string
}

type Users struct {
	db bob.DB
}

func (u *Users) Query() *sqlite.ViewQuery[*models.User, models.UserSlice] {
	return models.Users.Query(context.Background(), u.db)
}

func (u *Users) Create(s models.UserSetter) (*models.User, error) {
	return models.Users.Insert(context.Background(), u.db, &s)
}

func (u *Users) Exists(id int) (bool, error) {
	return models.UserExists(context.Background(), u.db, int32(id))
}

func (u *Users) Get(f UserF) (*models.User, error) {
	q := u.Query()
	if f.Role != "" {
		q.Apply(models.SelectWhere.Users.Role.EQ(string(f.Role)))
	}
	if f.Username != "" {
		q.Apply(models.SelectWhere.Users.Username.EQ(f.Username))
	}
	qb, _, _ := q.Build()
	slog.Debug("Get User", "query", qb)
	m, err := q.One()
	return m, catchErr("GET USER", err)
}

func catchErr(op string, err error) error {
	if err == nil {
		return nil
	}

	e := err.Error()
	slog.Error(op + ": " + e)
	switch e {
	case sql.ErrNoRows.Error():
		return fmt.Errorf("%w %s", ErrNotFound, op)
		// default:
		// 	if strings.Contains(err.Error(), "UNIQUE constraint failed:") {
		// 		return fmt.Errorf("%w %s", database.ErrUniqueCheckFailed, msg)
		// 	}
		// 	return errors.NewFatal(err.Error())
	}
	return err
}
