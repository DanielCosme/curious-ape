package repository

import (
	"github.com/danielcosme/curious-ape/internal/database"
	"testing"

	"github.com/danielcosme/curious-ape/internal/repository/sqlite"
	"github.com/jmoiron/sqlx"
)

func NewSqlite(db *sqlx.DB) *database.Repository {
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
	return NewSqlite(db)
}
