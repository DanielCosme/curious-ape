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
		Global.validate()
	})
}

type Config struct {
	Environment  Environment
	Port         string
	LogLevel     slog.Level
	DatabasePath string
	Username     string
	Password     string
}

func (c *Config) validate() {
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
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
