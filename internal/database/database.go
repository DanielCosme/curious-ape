package database

import (
	"errors"
	"github.com/stephenafamo/bob"
)

var (
	ErrNotFound           = errors.New("database: not found")
	ErrInvalidCredentials = errors.New("database: invalid credentials")
)

type Database struct {
	Users    Users
	Days     Days
	executor bob.DB
}

func New(executor bob.DB) *Database {
	return &Database{
		Users:    Users{db: executor},
		Days:     Days{db: executor},
		executor: executor,
	}
}
