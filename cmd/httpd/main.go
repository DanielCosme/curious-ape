package main

import (
	"log"
	"net/http"
	"time"

	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/internal/datasource/sqlite"
	"github.com/danielcosme/curious-ape/internal/transport"

	"github.com/jmoiron/sqlx"
)

type config struct {
	Database struct {
		DNS string
	}
}

func main() {
	// Load configuration: env, file, flags
	// Dial data sources: databases...
	// Initialize transport: server, router, logger, etc...
	//		Wrap and invoke application logic.
	// Listen and serve
	cfg := &config{}
	cfg.Database.DNS = "ape.db"

	db := sqlx.MustConnect(sqlite.DriverName, cfg.Database.DNS)

	api := &transport.Transport{
		App: application.New(db),
		Server: &http.Server{
			Addr:         ":4000",
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
	}

	if err := api.ListenAndServe(); err != nil {
		log.Fatal()
	}
}
