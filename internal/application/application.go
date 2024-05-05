package application

import (
	"github.com/danielcosme/curious-ape/internal/database"
	"github.com/danielcosme/curious-ape/internal/entity"
	"github.com/danielcosme/curious-ape/internal/integrations"
	"log/slog"
)

type App struct {
	Log  *slog.Logger
	db   *database.Repository
	cfg  *Environment
	sync *integrations.Sync
}

// Endpoints and application methods to sync manually
// 		Then I will set up cron-jobs on my linux server to invoke them, because they are hosted on the same machine
//		there is no need for authentication

type AppOptions struct {
	Logger     *slog.Logger
	Config     *Environment
	Repository *database.Repository
}

type Environment struct {
	Fitbit *entity.Oauth2Config
	Google *entity.Oauth2Config
	Env    string
}

func New(opts *AppOptions) *App {
	a := &App{
		Log:  opts.Logger,
		db:   opts.Repository,
		cfg:  opts.Config,
		sync: integrations.NewSync(opts.Logger),
	}

	a.Log.Info("Application initialized", "Environment", a.cfg.Env)
	return a
}

type props map[string]string
