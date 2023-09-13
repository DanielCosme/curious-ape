package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/repository"
	"github.com/danielcosme/curious-ape/internal/repository/sqlite"
	"github.com/danielcosme/curious-ape/internal/web"
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
		Port     int `json:"port"`
		FilePath string
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
	flag.StringVar(&cfg.Environment, "env", "", "Sets the running environment for the application")
	flag.Parse()
	setFilePath(cfg)
	readConfiguration(cfg)

	// logger initialization
	logger := logape.New(os.Stdout, logape.LevelDebug, time.RFC822)
	logape.DefaultLogger = logger

	// SQL datasource initialization
	db := sqlx.MustConnect(sqlite.DriverName, cfg.Server.FilePath+"/"+cfg.Database.DNS)

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db.DB)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode

	web := &web.WebClient{
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

	if err := web.App.SetPassword(cfg.Admin.Name, cfg.Admin.Password, entity.AdminRole); err != nil {
		logger.Fatal(err)
	}
	if err := web.App.SetPassword(cfg.User.Name, cfg.User.Password, entity.UserRole); err != nil {
		logger.Fatal(err)
	}
	if err := web.App.SetPassword(cfg.Guest.Name, cfg.Guest.Password, entity.GuestRole); err != nil {
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
