package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-co-op/gocron/v2"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/internal/core/entity"
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
		Fitbit *entity.Oauth2Config `json:"fitbit"`
		Google *entity.Oauth2Config `json:"google"`
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

func main() {
	// flags & configuration
	cfg := new(config)
	readConfiguration(cfg)

	// logger initialization
	logger := logape.New(os.Stdout, logape.LevelDebug, time.RFC822)
	logape.DefaultLogger = logger

	// SQL datasource initialization
	db := sqlx.MustConnect(sqlite.DriverName, "./"+cfg.Database.DNS)

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db.DB)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode

	web := &transport.WebClient{
		App: application.New(&application.AppOptions{
			Repository: repository.NewSqlite(db),
			Config: &application.Environment{
				Env:    cfg.Environment,
				Fitbit: cfg.Integrations.Fitbit,
				Google: cfg.Integrations.Google,
			},
			Logger:         logger,
			SessionManager: sessionManager,
		}),
		Server: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
			IdleTimeout:  time.Minute,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			ErrorLog:     log.New(logger, "", 0),
		},
	}

	go func() {
		logger.Info("Starting cron jobs")
		if err := startCron(web.App); err != nil {
			logger.Fatal(err)
		}
		logger.Info("Finished starting cron Jobs")
	}()

	if err := web.App.SetPassword(cfg.Admin.Name, cfg.Admin.Password, entity.AdminRole); err != nil {
		logger.Fatal(err)
	}
	if err := web.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}

func startCron(a *application.App) error {
	s, err := gocron.NewScheduler()
	if err != nil {
		return err
	}

	j, err := s.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(
			gocron.NewAtTime(23, 0, 0),
		)),
		gocron.NewTask(func() {
			err := a.SyncFitbitSleepLog(time.Now())
			if err != nil {
				a.Log.Error(fmt.Errorf("cron Job: %w", err))
			}
		}),
	)
	if err != nil {
		return err
	}
	s.Start()
	next, err := j.NextRun()
	if err != nil {
		a.Log.Fatal(err)
	}
	a.Log.Info("Cron job next run: ", next.Local())

	return err
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
