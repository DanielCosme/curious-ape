package application

import (
	"github.com/danielcosme/curious-ape/internal/database"
	"github.com/danielcosme/curious-ape/internal/integrations"
	"golang.org/x/oauth2"
	"log/slog"
)

type App struct {
	Log  *slog.Logger
	Env  Environment
	db   *database.Database
	sync *integrations.Integrations
}

// Endpoints and application methods to sync manually
// 		Then I will set up cron-jobs on my linux server to invoke them, because they are hosted on the same machine
//		there is no need for authentication

type AppOptions struct {
	Logger   *slog.Logger
	Config   *Config
	Database *database.Database
}

type Environment string

const (
	Dev     Environment = "dev"
	Test    Environment = "test"
	Prod    Environment = "prod"
	Staging Environment = "staging"
)

type Config struct {
	Fitbit *oauth2.Config
	Google *oauth2.Config
	Env    Environment
}

func New(opts *AppOptions) *App {
	s := integrations.New(opts.Config.Fitbit)
	a := &App{
		Log:  opts.Logger,
		Env:  opts.Config.Env,
		db:   opts.Database,
		sync: s,
	}

	a.Log.Info("Application initialized", "Config", a.Env)
	return a
}
