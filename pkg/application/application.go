package application

import (
	"log/slog"
)

type Application struct {
	// TODO: Maybe create en environment enum
	Env    string
	logger *slog.Logger
}

func New(logger *slog.Logger, env string) *Application {
	return &Application{
		logger: logger,
		Env:    env,
	}
}
