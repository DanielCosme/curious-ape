package repository

import (
	"testing"

	"github.com/danielcosme/curious-ape/internal/core/database"
	"github.com/danielcosme/curious-ape/internal/repository/sqlite"
	"github.com/jmoiron/sqlx"
)

func NewSqliteRepository(db *sqlx.DB) *database.Repository {
	return &database.Repository{
		Habits:      &sqlite.HabitsDataSource{DB: db},
		Days:        &sqlite.DaysDataSource{DB: db},
		Auths:       &sqlite.AuthenticationDataSource{DB: db},
		SleepLogs:   &sqlite.SleepLogDataSource{DB: db},
		FitnessLogs: &sqlite.FitnessLogDataSource{DB: db},
		Users:       &sqlite.UsersDataSource{DB: db},
	}
}

func NewTestSqliteRepository(t *testing.T) *database.Repository {
	t.Helper()
	db := sqlite.NewTestSqliteDB(t)
	return NewSqliteRepository(db)
}
