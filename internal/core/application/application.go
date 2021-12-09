package application

import (
	"database/sql"
	"github.com/danielcosme/curious-ape/internal/datasources/sqlite"
)

type App struct {
	Habits *HabitsInteractor
}

func New(db *sql.DB) *App {
	return &App{
		Habits: &HabitsInteractor{Service: sqlite.New(db)},
	}
}
