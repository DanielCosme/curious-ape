package main

import (
	"flag"
	"fmt"
	"github.com/danielcosme/curious-ape/pkg/application"
	"github.com/danielcosme/curious-ape/pkg/transport"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// TODO: Implement deployment to ape.danicos.dev (test environment).
// TODO: Implement deployment to ape.danicos.me (main environment).

const version = "1.0.0"

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "Server port")
	flag.StringVar(&cfg.env, "env", "development",
		"Operating environment (development|testing|staging|production)")
	flag.Parse()

	// TODO: Figure out options for debugging sessions, maybe I want AddSource=true in some cases
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	app := application.New(logger, cfg.env)
	t := transport.New(app, version)

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      transport.Router(t),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// TODO: handle graceful termination with server.Shutdown and server.Close.
	// 		These signals come from UNIX signals.
	slog.Info("Application started", "environment", cfg.env, "version", version)
	slog.Info("Server listening", "port", cfg.port)
	err := server.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

type config struct {
	port int
	env  string
}
