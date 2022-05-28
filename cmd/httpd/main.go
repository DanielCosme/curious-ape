package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/datasource/sqlite"
	"github.com/danielcosme/curious-ape/internal/transport"
	"github.com/danielcosme/curious-ape/sdk/logape"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Initialize configuration
	cfg := new(config)

	flag.StringVar(&cfg.Environment, "env", "", "Sets the running environment for the application")
	flag.Parse()

	// panics if there are any errors
	readConfiguration(cfg)
	db := sqlx.MustConnect(sqlite.DriverName, cfg.Database.DNS)
	logger := logape.New(os.Stdout, logape.LevelTrace, time.RubyDate)

	api := &transport.API{
		App: application.New(&application.AppOptions{
			DB: db,
			Config: &application.Environment{
				Env:    cfg.Environment,
				Fitbit: cfg.Integrations.Fitbit,
			},
			Logger: logger,
		}),
		Server: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
	}

	if err := api.Run(); err != nil {
		log.Fatal()
	}
}

type config struct {
	Database struct {
		DNS string `json:"dns"`
	} `json:"database"`
	Server struct {
		Port int `json:"port"`
	} `json:"server"`
	Integrations struct {
		Fitbit *entity.Oauth2Config `json:"fitbit"`
	} `json:"integrations"`
	Environment string `json:"environment"`
}

func readConfiguration(cfg *config) *config {
	var err error
	rawFile := []byte{}

	switch cfg.Environment {
	case "dev":
		rawFile, err = os.ReadFile(".dev.env.json")
	case "prod":
		rawFile, err = os.ReadFile(".env.json")
	default:
		err = errors.New("no environment provided")
	}
	panicIfErr(err)

	err = json.Unmarshal(rawFile, cfg)
	panicIfErr(err)
	return cfg
}

func panicIfErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
