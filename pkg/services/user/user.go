package user

import (
	"errors"
	"log/slog"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"github.com/stephenafamo/bob"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name     string
	Password string
}

type Service struct {
	db bob.DB
}

func NewService(db bob.DB) *Service {
	return &Service{db: db}
}

func (s *Service) SetPassword(username, password string) error {
	slog.Info("Setting password", "username", username)
	if password == "" {
		return errors.New("password cannot be empty")
	}
	if username == "" {
		return errors.New("username cannot be empty")
	}

	user, err := Get(s.db, Params{Username: username})
	if core.IfErrNNotFound(err) {
		return err
	}
	if user == nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			return err
		}
		if err := create(s.db, username, string(hash)); err != nil {
			return err
		}
	}
	return nil
}
