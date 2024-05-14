package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/danielcosme/curious-ape/internal/database"
	"github.com/lmittmann/tint"
	"github.com/stephenafamo/bob"
	"golang.org/x/oauth2"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"time"

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

	db, err := sql.Open("sqlite3", "./"+cfg.Database.DSN)
	exitIfErr(err)
	err = db.Ping()
	exitIfErr(err)

	app := application.New(&application.AppOptions{
		Database: database.New(bob.NewDB(db)),
		Config: &application.Config{
			Env:    cfg.Environment,
			Fitbit: cfg.Integrations.Fitbit.ToConf(),
			Google: cfg.Integrations.Google.ToConf(),
		},
		Logger: sLogger,
	})

	err = app.SetPassword(cfg.Admin.UserName, cfg.Admin.Password, cfg.Admin.Email, core.AdminRole)
	exitIfErr(err)

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	t, err := transport.NewTransport(app, sessionManager, v)
	exitIfErr(err)

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

	cfg.Environment = application.Environment(os.Getenv("APE_ENVIRONMENT"))
	if cfg.Environment == "" {
		logFatal(errors.New("environment variable APE_ENVIRONMENT is empty"))
	} else if cfg.Environment != application.Dev && cfg.Environment != application.Prod && application.Staging != cfg.Environment {
		logFatal(fmt.Errorf("invalid environment: '%s'", cfg.Environment))
	}
	rawFile, err = os.ReadFile("config.json")
	exitIfErr(err)

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
	dirty := false

	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				hash = s.Value[:7]
			case "vcs.modified":
				if s.Value == "true" {
					dirty = true
				}
			}
		}
	}
	if dirty {
		return fmt.Sprintf("%s-%s-dirty", version, hash)
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
