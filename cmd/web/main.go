package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/danielcosme/curious-ape/internal/database"
	"github.com/go-co-op/gocron/v2"
	"github.com/lmittmann/tint"
	"github.com/stephenafamo/bob"
	"golang.org/x/oauth2"

	"github.com/danielcosme/curious-ape/internal/application"
	"github.com/danielcosme/curious-ape/internal/core"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/danielcosme/curious-ape/internal/transport"
	_ "github.com/mattn/go-sqlite3"

	"github.com/alexedwards/scs/v2"
)

type config struct {
	Port     int `json:"port"`
	Database struct {
		DSN string `json:"dsn"`
	} `json:"database"`
	Integrations struct {
		Fitbit *Oauth2Config `json:"fitbit"`
		Google *Oauth2Config `json:"google"`
		Toggl  struct {
			Token       string `json:"api_token"`
			WorkspaceID int    `json:"workspace_id"`
		} `json:"toggl"`
	} `json:"integrations"`
	Environment application.Environment
	Admin       user `json:"admin"`
	User        user `json:"user"`
	Guest       user `json:"guest"`
}

type user struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

var version string

func main() {
	// flags & configuration
	cfg := new(config)
	v := Version()
	readConfiguration(cfg)

	// logger configuration
	logHandler := tint.NewHandler(os.Stdout, &tint.Options{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
		TimeFormat:  time.RFC822,
		NoColor:     false,
	})
	sLogger := slog.New(logHandler)
	slog.SetDefault(sLogger)

	sLogger.Info("Version: " + v)

	slog.Info("Opening database", "path", cfg.Database.DSN)
	db, err := sql.Open("sqlite3", cfg.Database.DSN)
	exitIfErr(err)
	err = db.Ping()
	exitIfErr(err)

	app := application.New(&application.AppOptions{
		Database: database.New(bob.NewDB(db)),
		Config: &application.Config{
			Env:              cfg.Environment,
			Fitbit:           cfg.Integrations.Fitbit.ToConf(),
			Google:           cfg.Integrations.Google.ToConf(),
			TogglToken:       cfg.Integrations.Toggl.Token,
			TogglWorkspaceID: cfg.Integrations.Toggl.WorkspaceID,
		},
		Logger: sLogger,
	})

	err = app.SetPassword(cfg.Admin.UserName, cfg.Admin.Password, cfg.Admin.Email, core.AuthRoleAdmin)
	exitIfErr(err)

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 48 * time.Hour
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	t, err := transport.NewTransport(app, sessionManager, v)
	exitIfErr(err)

	app.Log.Info("Launching cron jobs")
	if err := setUpCronJobs(app); err != nil {
		logFatal(err)
	}

	addr := fmt.Sprintf(":%d", cfg.Port)
	server := http.Server{
		Addr:         addr,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logHandler, slog.LevelError),
		Handler:      transport.EchoRoutes(t),
	}
	t.App.Log.Info("HTTP server listening", "addr", addr)
	if err := server.ListenAndServe(); err != nil {
		logFatal(err)
	}
}

func readConfiguration(cfg *config) *config {
	var err error
	var rawFile []byte

	env, err := application.ParseEnvironment(os.Getenv("APE_ENVIRONMENT"))
	if err != nil {
		logFatal(errors.New("environment variable APE_ENVIRONMENT is empty"))
	}
	cfg.Environment = env
	configPath := "config.json"
	rawFile, err = os.ReadFile(configPath)
	exitIfErr(err)
	slog.Info("Configuration file loaded", "path", configPath)

	err = json.Unmarshal(rawFile, cfg)
	exitIfErr(err)
	return cfg
}

func exitIfErr(err error) {
	if err != nil {
		logFatal(err)
	}
}

func logFatal(err error) {
	slog.Error("Fatal failure", "err", err.Error(), "stack", string(debug.Stack()))
	os.Exit(1)
}

func Version() string {
	hash := "unknown"
	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				hash = s.Value[:7]
			}
		}
	}
	return fmt.Sprintf("%s-%s", version, hash)
}

type Oauth2Config struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURL  string   `json:"redirect_url"`
	TokenURL     string   `json:"token_url"`
	AuthURL      string   `json:"auth_url"`
	AuthStyle    int      `json:"auth_style"`
	Scopes       []string `json:"scopes"`
}

func (o Oauth2Config) ToConf() *oauth2.Config {
	slog.Info("Loading Oauth2 configuration", "redirect", o.RedirectURL)
	return &oauth2.Config{
		ClientID:     o.ClientID,
		ClientSecret: o.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   o.AuthURL,
			TokenURL:  o.TokenURL,
			AuthStyle: oauth2.AuthStyle(o.AuthStyle), // Zero value means auto-detect.
		},
		RedirectURL: o.RedirectURL,
		Scopes:      o.Scopes,
	}
}

func setUpCronJobs(a *application.App) error {
	s, err := gocron.NewScheduler()
	if err != nil {
		return err
	}

	today := core.NewDateToday()
	_, err = s.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(
			gocron.NewAtTime(23, 0, 0),
		)),
		gocron.NewTask(func() {
			if err := a.SleepSync(today); err != nil {
				a.Log.Error(err.Error())
			}
		}),
		gocron.WithName("Sleep logs sync"),
	)
	if err != nil {
		return err
	}
	_, err = s.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(
			gocron.NewAtTime(23, 55, 0),
		)),
		gocron.NewTask(func() {
			if err := a.DeepWorkSync(today); err != nil {
				a.Log.Error(err.Error())
			}
		}),
		gocron.WithName("Deep work logs sync"),
	)
	if err != nil {
		return err
	}

	_, err = s.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(
			gocron.NewAtTime(23, 55, 0),
		)),
		gocron.NewTask(func() {
			if err := a.FitnessSync(today); err != nil {
				a.Log.Error(err.Error())
			}
		}),
		gocron.WithName("Fitness logs sync"),
	)
	if err != nil {
		return err
	}
	s.Start()

	for _, j := range s.Jobs() {
		next, err := j.NextRun()
		if err != nil {
			return err
		}
		a.Log.Info("Cron job configured", "name", j.Name(), "next_run", next.Format(core.HumanDateWithTime), "Timezone", next.Location().String())
	}

	return nil
}
