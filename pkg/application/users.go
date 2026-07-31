package application

import (
	"danicos.dev/daniel/curious-ape/pkg/gen/bob/models"
	"danicos.dev/daniel/curious-ape/pkg/persistence"
)

func (a *App) UserExists(id int) (bool, error) {
	return a.db.Users.Exists(id)
}

func (a *App) GetUser(id int) (*models.User, error) {
	return a.db.Users.Get(persistence.UserParams{ID: id})
}
