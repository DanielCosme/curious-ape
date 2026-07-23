package config

import (
	"log/slog"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Environment string

const (
	Dev  Environment = "dev"
	CI   Environment = "ci"
	Prod Environment = "prod"
)

var (
	Global *Config
	once   sync.Once
)

func init() {
	once.Do(func() {
		Global = Load()
	})
}

type Config struct {
	Environment  Environment
	Port         string
	LogLevel     slog.Level
	DatabasePath string
}

func loadBase() *Config {
	godotenv.Load()

	return &Config{
		Port:         getEnv(CONFIG_PORT, "4000"),
		DatabasePath: getEnv(CONFIG_DB_PATH, TMP_DIR+"/ape.db"),
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
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
