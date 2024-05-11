package database

import (
	"github.com/stephenafamo/bob"
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
