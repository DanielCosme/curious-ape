package datasource

import (
	"github.com/danielcosme/curious-ape/internal/core/repository"
	"github.com/danielcosme/curious-ape/internal/datasource/sqlite"
	"github.com/jmoiron/sqlx"
)

type DataModel struct {
	Habits repository.Habit
}

func New(db *sqlx.DB) DataModel {
	return DataModel{
		Habits: sqlite.NewHabitsDataSource(db),
	}
}
