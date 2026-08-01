package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Environment string

const (
	Dev  Environment = "dev"
	CI   Environment = "ci"
	Prod Environment = "prod"
)

type Config struct {
	Environment      Environment
	Port             string
	LogLevel         slog.Level
	DatabasePath     string
	Username         string
	Password         string
	HevyAPIKey       string
	TogglAPIKey      string
	TogglWorkspaceID int
}

func (c *Config) Validate() {
	if c.DatabasePath == "" {
		panic("config error: database path empty")
	}
	if c.Username == "" {
		panic("config error: username empty")
	}
	if c.Password == "" {
		panic("config error: password empty")
	}
}

func IsDev() bool {
	if value, _ := os.LookupEnv(ENVIRONMENT); value == string(Dev) {
		return true
	}
	return false
}

func loadBase() *Config {
	godotenv.Load()

	return &Config{
		Port: getEnv(CONFIG_PORT, "4000"),
		LogLevel: func() slog.Level {
			switch os.Getenv(CONFIG_LOG_LEVEL) {
			case slog.LevelDebug.String():
				return slog.LevelDebug
			case slog.LevelInfo.String():
				return slog.LevelInfo
			case slog.LevelWarn.String():
				return slog.LevelWarn
			case slog.LevelError.String():
				return slog.LevelError
			default:
				return slog.LevelInfo
			}
		}(),
		HevyAPIKey:  getEnv(CONFIG_HEVY_API_KEY, ""),
		TogglAPIKey: getEnv(CONFIG_TOGGL_API_KEY, ""),
		TogglWorkspaceID: func() int {
			wokID := getEnv(CONFIG_TOGGL_WORKSPACE_ID, "")
			id, err := strconv.Atoi(wokID)
			if err != nil {
				panic(err)
			}
			return id
		}(),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		slog.Info(fmt.Sprintf("Config value loaded from environment: %s", key))
		return val
	}
	return fallback
}
