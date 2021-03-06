package application

import (
	"github.com/danielcosme/curious-ape/internal/core/database"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/integrations"
	"github.com/danielcosme/curious-ape/internal/repository/sqlite"
	"github.com/danielcosme/curious-ape/sdk/log"

	"github.com/jmoiron/sqlx"
)

type App struct {
	db   *database.Models
	cfg  *Environment
	Log  *log.Logger
	Sync *integrations.Sync
}

// Endpoints and application methods to sync manually
// 		Then I will set up cron-jobs on my linux server to invoke them, because they are hosted on the same machine
//		there is no need for authentication

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
		db: &database.Models{
			Habits:    &sqlite.HabitsDataSource{DB: opts.DB},
			Days:      &sqlite.DaysDataSource{DB: opts.DB},
			Oauths:    &sqlite.Oauth2DataSource{DB: opts.DB},
			SleepLogs: &sqlite.SleepLogDataSource{DB: opts.DB},
		},
		cfg:  opts.Config,
		Log:  opts.Logger,
		Sync: integrations.NewSync(opts.Logger),
	}

	a.Log.InfoP("Application running", log.Prop{"environment": a.cfg.Env})
	return a
}
