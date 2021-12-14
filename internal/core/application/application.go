package application

import (
	"database/sql"

	"github.com/danielcosme/curious-ape/internal/core/repository"
	"github.com/danielcosme/curious-ape/internal/datasource/sqlite"
)

type App struct {
	Habits repository.Habit
}

func New(db *sql.DB) *App {
	return &App{
		Habits: sqlite.NewHabitsService(db),
	}
}

func (a *App) tester() {
}
