package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/internal/datasource/sqlite"
	"github.com/danielcosme/curious-ape/internal/transport/httprest"

	_ "github.com/mattn/go-sqlite3"
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

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	api := &httprest.API{
		App: application.New(db),
	}
	
	api.Server = &http.Server{
		Addr:              ":4000",
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if err := api.Run() ; err != nil {
		log.Fatal()
	}
}

func openDB(cfg *config) (*sql.DB, error) {
	db, err := sql.Open(sqlite.DriverName, cfg.Database.DNS)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
