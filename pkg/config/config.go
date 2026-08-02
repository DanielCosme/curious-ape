package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

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
	Fitbit           OathConfig
}

type OathConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	TokenURL     string
	AuthURL      string
	Scopes       []string
	AuthStyle    int
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
		Port:             getEnv(CONFIG_PORT, "4000"),
		HevyAPIKey:       getEnv(CONFIG_HEVY_API_KEY, ""),
		TogglAPIKey:      getEnv(CONFIG_TOGGL_API_KEY, ""),
		TogglWorkspaceID: parseInt(getEnv(CONFIG_TOGGL_WORKSPACE_ID, "")),
		Fitbit: OathConfig{
			ClientID:     getEnv(CONFIG_FITBIT_CLIENT_ID, ""),
			ClientSecret: getEnv(CONFIG_FITBIT_CLIENT_SECRET, ""),
			RedirectURL:  getEnv(CONFIG_FITBIT_REDIRECT_URL, ""),
			TokenURL:     getEnv(CONFIG_FITBIT_TOKEN_URL, ""),
			AuthURL:      getEnv(CONFIG_FITBIT_AUTH_URL, ""),
			Scopes:       parseArray(getEnv(CONFIG_FITBIT_SCOPES, "")),
		},
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
		slog.Info(fmt.Sprintf("Config value loaded from environment: %s", key))
		return val
	}
	return fallback
}

func parseArray(values string) []string {
	return strings.Split(values, ",")
}

func parseInt(value string) int {
	id, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return id
}
