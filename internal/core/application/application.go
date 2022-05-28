package application

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/core/repository"
	"github.com/danielcosme/curious-ape/internal/datasource/sqlite"
	"github.com/danielcosme/curious-ape/sdk/logape"
	"github.com/jmoiron/sqlx"
)

type App struct {
	db  *repository.Models
	cfg *Environment
	log *logape.Logger
}

type AppOptions struct {
	DB     *sqlx.DB
	Logger *logape.Logger
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
		log: opts.Logger,
	}

	props := map[string]string{}
	props["environment"] = a.cfg.Env

	a.log.Info("Applications running in ", a.cfg.Env , " environment")
	a.log.InfoP("Application running", props)
	a.log.Infof("Applications running in %s environment", a.cfg.Env)
	return a
}

func (a *App) Error(err error, properties map[string]string) {
	a.log.ErrorP(err, properties)
}
