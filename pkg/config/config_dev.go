//go:build dev || !prod

package config

import (
	"os"
)

func Load() *Config {
	cfg := loadBase()
	cfg.Environment = Dev
	cfg.DatabasePath = getEnv(CONFIG_DB_PATH, TMP_DIR+"/ape.db")
	cfg.Username = getEnv(CONFIG_USERNAME, "admin")
	cfg.Password = getEnv(CONFIG_PASSWORD, "admin")

	os.Setenv(ENVIRONMENT, string(cfg.Environment))
	return cfg
}
