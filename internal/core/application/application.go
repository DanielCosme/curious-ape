package application

import (
	"github.com/alexedwards/scs/v2"
	"github.com/danielcosme/curious-ape/internal/core/database"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/integrations"
	"github.com/danielcosme/go-sdk/log"
)

type App struct {
	db   *database.Repository
	cfg  *Environment
	Log  *log.Logger
	sync *integrations.Sync
	// TODO(daniel) move the session manager to the transport layer.
	Session *scs.SessionManager
}

// Endpoints and application methods to sync manually
// 		Then I will set up cron-jobs on my linux server to invoke them, because they are hosted on the same machine
//		there is no need for authentication

type AppOptions struct {
	Logger         *log.Logger
	Config         *Environment
	Repository     *database.Repository
	SessionManager *scs.SessionManager
}

type Environment struct {
	Fitbit *entity.Oauth2Config
	Google *entity.Oauth2Config
	Env    string
}

func New(opts *AppOptions) *App {
	a := &App{
		db:      opts.Repository,
		cfg:     opts.Config,
		Log:     opts.Logger,
		sync:    integrations.NewSync(opts.Logger),
		Session: opts.SessionManager,
	}

	a.Log.InfoP("Application running", log.Prop{"environment": a.cfg.Env})
	return a
}

type props map[string]string
