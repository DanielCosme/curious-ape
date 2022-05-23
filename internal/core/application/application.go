package application

import (
	"github.com/danielcosme/curious-ape/internal/core/repository"
	"github.com/danielcosme/curious-ape/internal/datasource/sqlite"
	"github.com/jmoiron/sqlx"
)

type App struct {
	db *repository.Models
}

func New(db *sqlx.DB) *App {
	return &App{
		db: &repository.Models{
			Habits: &sqlite.HabitsDataSource{DB: db},
			Days:   &sqlite.DaysDataSource{DB: db},
		},
	}
}
