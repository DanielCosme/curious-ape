package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/danielcosme/curious-ape/internal/database"
	"github.com/lmittmann/tint"
	"github.com/stephenafamo/bob"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/danielcosme/curious-ape/internal/application"
	"github.com/danielcosme/curious-ape/internal/core"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/danielcosme/curious-ape/internal/transport"
	_ "github.com/mattn/go-sqlite3"

	"github.com/alexedwards/scs/v2"
)

type config struct {
	Database struct {
		DSN string `json:"dns"`
	} `json:"database"`
	Server struct {
		Port int `json:"port"`
	} `json:"server"`
	Integrations struct {
		Fitbit *core.Oauth2Config `json:"fitbit"`
		Google *core.Oauth2Config `json:"google"`
	} `json:"integrations"`
	Environment string `json:"environment"`
	Admin       user   `json:"admin"`
	User        user   `json:"user"`
	Guest       user   `json:"guest"`
}

type user struct {
	Name     string `json:"name"`
	Password string `json:"password"`
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
			Fitbit: cfg.Integrations.Fitbit,
			Google: cfg.Integrations.Google,
		},
		Logger: sLogger,
	})

	err = app.SetPassword("daniel", "test", core.AdminRole)
	exitIfErr(err)

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	t, err := transport.NewTransport(app, sessionManager, v)
	exitIfErr(err)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
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

	cfg.Server.Port, err = strconv.Atoi(os.Getenv("APE_PORT"))
	if err != nil {
		logFatal(fmt.Errorf("invalid APE_PORT: '%w'", err))
	}

	cfg.Environment = os.Getenv("APE_ENVIRONMENT")
	if cfg.Environment == "" {
		logFatal(errors.New("environment variable APE_ENVIRONMENT is empty"))
	} else if cfg.Environment != "dev" && cfg.Environment != "prod" {
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
				hash = s.Value[:8]
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
