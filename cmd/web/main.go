package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/danielcosme/curious-ape/internal/application"
	entity2 "github.com/danielcosme/curious-ape/internal/entity"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/danielcosme/curious-ape/internal/repository"
	"github.com/danielcosme/curious-ape/internal/repository/sqlite"
	"github.com/danielcosme/curious-ape/internal/transport"
	"github.com/danielcosme/go-sdk/errors"
	logape "github.com/danielcosme/go-sdk/log"
	_ "github.com/mattn/go-sqlite3"

	"github.com/alexedwards/scs/v2"
	"github.com/jmoiron/sqlx"
)

type config struct {
	Database struct {
		DNS string `json:"dns"`
	} `json:"database"`
	Server struct {
		Port int `json:"port"`
	} `json:"server"`
	Integrations struct {
		Fitbit *entity2.Oauth2Config `json:"fitbit"`
		Google *entity2.Oauth2Config `json:"google"`
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

	// logger initialization
	logger := logape.New(os.Stdout, logape.LevelDebug, time.RFC822)
	logape.DefaultLogger = logger
	logger.Info("Version ", v)

	// SQL datasource initialization
	db := sqlx.MustConnect(sqlite.DriverName, "./"+cfg.Database.DNS)

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db.DB)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode

	app := application.New(&application.AppOptions{
		Repository: repository.NewSqlite(db),
		Config: &application.Environment{
			Env:    cfg.Environment,
			Fitbit: cfg.Integrations.Fitbit,
			Google: cfg.Integrations.Google,
		},
		Logger: logger,
	})

	t, err := transport.NewTransport(app, v, sessionManager)
	if err != nil {
		exitIfErr(err)
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	server := http.Server{
		Addr:         addr,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     log.New(logger, "", 0),
		Handler:      transport.EchoRoutes(t),
	}
	t.App.Log.InfoP("HTTP server listening", logape.Prop{"addr": addr})
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}

func readConfiguration(cfg *config) *config {
	var err error
	var rawFile []byte
	cfg.Environment = os.Getenv("APE_ENVIRONMENT")
	cfg.Server.Port, err = strconv.Atoi(os.Getenv("APE_PORT"))
	if err != nil {
		logape.DefaultLogger.Fatal(fmt.Errorf("Invalid APE_PORT: '%w'", err))
	}

	if cfg.Environment != "dev" && cfg.Environment != "prod" {
		if cfg.Environment == "" {
			logape.DefaultLogger.Fatal(errors.NewFatal("Environment variable APE_ENVIRONMENT is empty"))
		} else {
			logape.DefaultLogger.Fatal(errors.NewFatal(fmt.Sprintf("Invalid environment: '%s'", cfg.Environment)))
		}
	}
	rawFile, err = os.ReadFile("config.json")
	exitIfErr(err)

	err = json.Unmarshal(rawFile, cfg)
	exitIfErr(err)
	return cfg
}

func exitIfErr(err error) {
	if err != nil {
		logape.DefaultLogger.Fatal(err)
	}
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
