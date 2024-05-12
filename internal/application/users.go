package application

import (
	"errors"
	"github.com/aarondl/opt/omit"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database"
	"github.com/danielcosme/curious-ape/internal/database/gen/models"

	"golang.org/x/crypto/bcrypt"
)

func (a *App) SetPassword(username, password, email string, role core.AuthRole) error {
	a.Log.Info("Setting password", "username", username, "role", role)
	if password == "" {
		return errors.New("password cannot be empty")
	}
	if username == "" {
		return errors.New("username cannot be empty")
	}

	u, err := a.db.Users.Get(database.UserF{Role: role, Username: username})
	if database.IfNotFoundErr(err) {
		return err
	}
	if u == nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			return err
		}
		_, err = a.db.Users.Create(models.UserSetter{
			Username: omit.From(username),
			Password: omit.From(string(hash)),
			Role:     omit.From(string(role)),
			Email:    omit.From(email),
		})
		return err
	}
	return nil
}

// Authenticate returns userID if successfully authenticated.
func (a *App) Authenticate(username, password string) (int, error) {
	u, err := a.db.Users.Get(database.UserF{Username: username})
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return 0, database.ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, database.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return int(u.ID), nil
}

func (a *App) UserExists(id int) (bool, error) {
	return a.db.Users.Exists(id)
}
