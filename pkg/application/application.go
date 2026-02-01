package application

import (
	"errors"
	"log/slog"

	"github.com/danielcosme/curious-ape/pkg/integrations"
	"github.com/danielcosme/curious-ape/pkg/oak"
	"github.com/danielcosme/curious-ape/pkg/persistence"
	"golang.org/x/oauth2"
)

type App struct {
	Log  *oak.Oak // Maybe delete the logger.
	Env  Environment
	db   *persistence.Database
	sync *integrations.Integrations
}

type AppOptions struct {
	Logger   *oak.Oak
	Config   *Config
	Database *persistence.Database
}

type Config struct {
	Fitbit           *oauth2.Config
	Google           *oauth2.Config
	TogglToken       string
	TogglWorkspaceID int
	HevyAPIKey       string
	Env              Environment
}

func New(opts *AppOptions) *App {
	sync, _ := integrations.New(opts.Config.TogglWorkspaceID, opts.Config.TogglToken, opts.Config.HevyAPIKey, opts.Config.Fitbit, opts.Config.Google)
	a := &App{
		Log:  opts.Logger.Layer("app"),
		Env:  opts.Config.Env,
		db:   opts.Database,
		sync: sync,
	}
	a.Log.Info("Application initialized", "Environment", a.Env)
	return a
}

type Environment string

const (
	Prod Environment = "prod"
	Dev  Environment = "dev"
	Test Environment = "test"
)

func ParseEnvironment(s string) (Environment, error) {
	switch Environment(s) {
	case Prod:
		return Prod, nil
	case Dev:
		return Dev, nil
	case Test:
	case "":
		e := errors.New("empty environment field")
		slog.Error(e.Error())
		return "", e
	}
	e := errors.New("ivalid environment value: " + s)
	slog.Error(e.Error())
	return "", e
}
