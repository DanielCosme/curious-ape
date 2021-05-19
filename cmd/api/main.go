package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// It will later be created automatically during the build process
const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string // data source name
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

// http handler dependencies.
type application struct {
	config config
	logger *log.Logger
	debug  *log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 3000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Running environment")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("APE_PG_DB_DSN"),
		"PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25,
		"PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25,
		"PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m",
		"PostgreSQL max connection idle time")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	debug := log.New(os.Stdout, "DEBUG =>", log.Ldate|log.Ltime)

	// initialize db connection pool
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Printf("database connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
		debug:  debug,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting %s server on %d", cfg.env, cfg.port)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
