package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/danielcosme/curious-ape/internal/api"
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/repository"
	"github.com/danielcosme/curious-ape/internal/repository/sqlite"
	"github.com/danielcosme/go-sdk/errors"
	logape "github.com/danielcosme/go-sdk/log"

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

	api := &api.Transport{
		App: application.New(&application.AppOptions{
			Repository: repository.NewSqlite(db),
			Config: &application.Environment{
				Env:    cfg.Environment,
				Fitbit: cfg.Integrations.Fitbit,
				Google: cfg.Integrations.Google,
			},
			Logger: logger,
		}),
		Server: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
			ErrorLog:     log.New(logger, "", 0),
		},
	}

	if err := api.Run(); err != nil {
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
