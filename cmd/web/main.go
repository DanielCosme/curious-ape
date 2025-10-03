package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/stephenafamo/bob"
	"golang.org/x/oauth2"

	"github.com/danielcosme/curious-ape/pkg/application"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/oak"
	"github.com/danielcosme/curious-ape/pkg/persistence"

	"github.com/danielcosme/curious-ape/pkg/api"
	_ "modernc.org/sqlite"
)

var version string

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

func main() {
	logger := oak.New(oak.TintHandler(os.Stdout, oak.LevelTrace))
	oak.SetDefault(logger)
	logLogger := oak.NewLogLogger(logger, oak.LevelTrace)

	cfg := new(config)
	v := Version()
	readConfiguration(cfg)

	oak.Info("Version: " + v)
	oak.Info("Opening database", "path", cfg.Database.DSN)
	db, err := sql.Open("sqlite", cfg.Database.DSN)
	exitIfErr(err)
	err = db.Ping()
	exitIfErr(err)

	app := application.New(&application.AppOptions{
		Database: persistence.New(bob.NewDB(db)),
		Config: &application.Config{
			Env:              cfg.Environment,
			Fitbit:           cfg.Integrations.Fitbit.ToConf(),
			Google:           cfg.Integrations.Google.ToConf(),
			TogglToken:       cfg.Integrations.Toggl.Token,
			TogglWorkspaceID: cfg.Integrations.Toggl.WorkspaceID,
		},
		Logger: logger,
	})

	err = app.SetPassword(cfg.Admin.UserName, cfg.Admin.Password, cfg.Admin.Email, core.AuthRoleAdmin)
	exitIfErr(err)

	t := api.NewApi(app, v)

	addr := fmt.Sprintf(":%d", cfg.Port)
	server := http.Server{
		Addr:         addr,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     logLogger,
		Handler:      api.Routes(t),
	}
	oak.Info("HTTP server listening", "addr", addr)
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
	oak.Info("Configuration file loaded", "path", configPath)

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
	oak.Fatal("Fatal failure", "err", err.Error(), "stack", string(debug.Stack()))
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
	oak.Info("Loading Oauth2 configuration", "redirect", o.RedirectURL)
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
