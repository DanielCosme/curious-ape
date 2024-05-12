package application

import (
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database"
	"github.com/danielcosme/curious-ape/internal/integrations"
	"log/slog"
)

type App struct {
	Log  *slog.Logger
	db   *database.Database
	cfg  *Config
	sync *integrations.Sync
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
	Fitbit *core.Oauth2Config
	Google *core.Oauth2Config
	Env    Environment
}

func New(opts *AppOptions) *App {
	a := &App{
		Log:  opts.Logger,
		db:   opts.Database,
		cfg:  opts.Config,
		sync: integrations.NewSync(opts.Logger),
	}

	a.Log.Info("Application initialized", "Config", a.cfg.Env)
	return a
}
