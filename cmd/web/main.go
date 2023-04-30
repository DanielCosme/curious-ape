package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/repository/sqlite"
	"github.com/danielcosme/curious-ape/internal/web"
	"github.com/danielcosme/go-sdk/errors"
	logape "github.com/danielcosme/go-sdk/log"
	_ "github.com/mattn/go-sqlite3"

	"github.com/jmoiron/sqlx"
)

type config struct {
	Database struct {
		DNS string `json:"dns"`
	} `json:"database"`
	Server struct {
		Port     int `json:"port"`
		FilePath string
	} `json:"server"`
	Integrations struct {
		Fitbit *entity.Oauth2Config `json:"fitbit"`
		Google *entity.Oauth2Config `json:"google"`
	} `json:"integrations"`
	Environment string `json:"environment"`
	Password    string `json:"password"`
}

func main() {
	// flags & configuration
	cfg := new(config)
	flag.StringVar(&cfg.Environment, "env", "", "Sets the running environment for the application")
	flag.Parse()
	setFilePath(cfg)
	readConfiguration(cfg)

	// logger initialization
	logger := logape.New(os.Stdout, logape.LevelTrace, time.RFC822)
	logape.DefaultLogger = logger

	// SQL datasource initialization
	db := sqlx.MustConnect(sqlite.DriverName, cfg.Server.FilePath+"/"+cfg.Database.DNS)
	dbOther, err := sql.Open(sqlite.DriverName, cfg.Server.FilePath+"/"+cfg.Database.DNS)
	if err != nil {
		logger.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(dbOther)
	sessionManager.Lifetime = 12 * time.Hour

	web := &web.WebClient{
		App: application.New(&application.AppOptions{
			DB: db,
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

	if err := web.App.SetPassword(cfg.Password); err != nil {
		logger.Fatal(err)
	}

	if err := web.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}

func setFilePath(cfg *config) {
	path := fmt.Sprintf("%s/.ape/server", os.Getenv("HOME"))
	if err := os.MkdirAll(path, os.ModePerm); err != nil { // $HOME/.ape/server
		logape.DefaultLogger.Fatal(err)
	}

	cfg.Server.FilePath = path
}

func readConfiguration(cfg *config) *config {
	var err error
	rawFile := []byte{}
	filePath := cfg.Server.FilePath + "/"

	switch cfg.Environment {
	case "dev":
		filePath = filePath + "dev.env.json"
	case "prod":
		filePath = filePath + "prod.env.json"
	default:
		logape.DefaultLogger.Fatal(errors.NewFatal("no valid environment provided"))
	}
	rawFile, err = os.ReadFile(filePath)
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
