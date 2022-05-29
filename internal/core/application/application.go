package application

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/core/repository"
	"github.com/danielcosme/curious-ape/internal/datasource/sqlite"
	"github.com/danielcosme/curious-ape/sdk/log"
	"github.com/jmoiron/sqlx"
)

type App struct {
	db  *repository.Models
	cfg *Environment
	Log *log.Logger
}

type AppOptions struct {
	DB     *sqlx.DB
	Logger *log.Logger
	Config *Environment
}

type Environment struct {
	Fitbit *entity.Oauth2Config
	Env    string
}

func New(opts *AppOptions) *App {
	a := &App{
		db: &repository.Models{
			Habits: &sqlite.HabitsDataSource{DB: opts.DB},
			Days:   &sqlite.DaysDataSource{DB: opts.DB},
			Oauths: &sqlite.Oauth2DataSource{DB: opts.DB},
		},
		cfg: opts.Config,
		Log: opts.Logger,
	}

	a.Log.InfoP("Application running", log.Prop{"environment": a.cfg.Env})
	return a
}
