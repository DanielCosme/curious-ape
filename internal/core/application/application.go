package application

import (
	"github.com/danielcosme/curious-ape/internal/datasource"
	"github.com/danielcosme/curious-ape/internal/datasource/sqlite"
	"github.com/jmoiron/sqlx"
)

type App struct {
	db     datasource.DataModel
	Habits HabitsInteractor
}

func New(db *sqlx.DB) *App {
	return &App{
		db:     datasource.New(db),
		Habits: HabitsInteractor{sqlite.HabitsDataSource{DB: db}},
	}
}
