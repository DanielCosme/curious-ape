package main

import (
	"encoding/json"
	"fmt"
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/internal/datasource/sqlite"
	"github.com/danielcosme/curious-ape/internal/transport"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// panics if there are any errors
	cfg := readConfiguration()
	db := sqlx.MustConnect(sqlite.DriverName, cfg.Database.DNS)

	api := &transport.Transport{
		App: application.New(&application.AppOptions{
			DB: db,
		}),
		Server: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
	}

	if err := api.ListenAndServe(); err != nil {
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
}

func readConfiguration() *config {
	cfg := new(config)
	file, err := os.ReadFile(".env.json")
	panicIfErr(err)

	err = json.Unmarshal(file, cfg)
	panicIfErr(err)
	return cfg
}

func panicIfErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
