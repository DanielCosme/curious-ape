package application

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/core/repository"
	"github.com/danielcosme/curious-ape/internal/datasource/sqlite"
	"github.com/jmoiron/sqlx"
)

type App struct {
	db  *repository.Models
	env *Environment
}

type AppOptions struct {
	ENV *Environment
	DB  *sqlx.DB
}

type Environment struct {
	Fitbit *entity.Oauth2Config
}

func New(opts *AppOptions) *App {
	return &App{
		db: &repository.Models{
			Habits: &sqlite.HabitsDataSource{DB: opts.DB},
			Days:   &sqlite.DaysDataSource{DB: opts.DB},
			Oauths: &sqlite.Oauth2DataSource{DB: opts.DB},
		},
		env: opts.ENV,
	}
}
