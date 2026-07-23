package application

import (
	"danicos.dev/daniel/curious-ape/pkg/apps/day"
	"danicos.dev/daniel/curious-ape/pkg/apps/habit"
	"danicos.dev/daniel/curious-ape/pkg/config"
	"danicos.dev/daniel/curious-ape/pkg/event"
	"danicos.dev/daniel/curious-ape/pkg/integrations"
	"danicos.dev/daniel/curious-ape/pkg/oak"
	"danicos.dev/daniel/curious-ape/pkg/persistence"
	"golang.org/x/oauth2"
)

type App struct {
	Log   *oak.Oak // Maybe delete the logger.
	Env   config.Environment
	db    *persistence.Database
	sync  *integrations.Integrations
	Day   *day.App
	Habit *habit.App
	Bus   event.Bus
}

type AppOptions struct {
	Logger   *oak.Oak
	Config   *Config
	Database *persistence.Database
	Day      *day.App
	Habit    *habit.App
	Bus      event.Bus
}

type Config struct {
	Fitbit           *oauth2.Config
	Google           *oauth2.Config
	TogglToken       string
	TogglWorkspaceID int
	HevyAPIKey       string
	Env              config.Environment
}

func New(opts *AppOptions) *App {
	sync := integrations.New(opts.Config.TogglWorkspaceID, opts.Config.TogglToken, opts.Config.HevyAPIKey, opts.Config.Fitbit, opts.Config.Google)
	a := &App{
		Log:   opts.Logger.Layer("app"),
		Env:   opts.Config.Env,
		db:    opts.Database,
		sync:  sync,
		Day:   opts.Day,
		Habit: opts.Habit,
		Bus:   opts.Bus,
	}
	a.Log.Info("Application initialized", "Environment", a.Env)
	return a
}
