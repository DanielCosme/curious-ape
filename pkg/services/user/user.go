package user

import (
	"errors"
	"log/slog"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/persistence"
	"github.com/stephenafamo/bob"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	core.RepositoryID
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

func (s *Service) Authenticate(username, password string) (id int64, err error) {
	user, err := Get(s.db, Params{Username: username})
	if err != nil {
		if errors.Is(err, core.ErrRepositoryNotFound) {
			return 0, persistence.ErrInvalidCredentials
		}
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, persistence.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return user.ID, nil
}
