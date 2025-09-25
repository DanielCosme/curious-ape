package application

import (
	"errors"
	"github.com/aarondl/opt/omit"
	"github.com/danielcosme/curious-ape/database/gen/models"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/persistence"

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

	u, err := a.db.Users.Get(persistence.UserParams{Role: role, Username: username})
	if core.IfErrNNotFound(err) {
		return err
	}
	if u == nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			return err
		}
		_, err = a.db.Users.Create(&models.UserSetter{
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
	u, err := a.db.Users.Get(persistence.UserParams{Username: username})
	if err != nil {
		if errors.Is(err, persistence.ErrNotFound) {
			return 0, persistence.ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, persistence.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return int(u.ID), nil
}

func (a *App) UserExists(id int) (bool, error) {
	return a.db.Users.Exists(id)
}

func (a *App) GetUser(id int) (*models.User, error) {
	return a.db.Users.Get(persistence.UserParams{ID: id})
}
