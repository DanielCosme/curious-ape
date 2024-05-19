package application

import (
	"errors"
	"log/slog"

	"github.com/danielcosme/curious-ape/internal/database"
	"github.com/danielcosme/curious-ape/internal/integrations"
	"golang.org/x/oauth2"
)

type App struct {
	Log  *slog.Logger
	Env  Environment
	db   *database.Database
	sync *integrations.Integrations
}

type AppOptions struct {
	Logger   *slog.Logger
	Config   *Config
	Database *database.Database
}

type Config struct {
	Fitbit *oauth2.Config
	Google *oauth2.Config
	Env    Environment
}

func New(opts *AppOptions) *App {
	a := &App{
		Log:  opts.Logger,
		Env:  opts.Config.Env,
		db:   opts.Database,
		sync: integrations.New(opts.Config.Fitbit),
	}
	a.Log.Info("Application initialized", "Environment", a.Env)
	return a
}

type Environment string

const (
	Prod    Environment = "prod"
	Dev     Environment = "dev"
	Test    Environment = "test"
	Staging Environment = "staging"
)

func ParseEnvironment(s string) (Environment, error) {
	switch Environment(s) {
	case Prod:
		return Prod, nil
	case Dev:
		return Dev, nil
	case Staging:
		return Staging, nil
	case Test:
	}
	return "", errors.New("Invalid env value: " + s)
}
