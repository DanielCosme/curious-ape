//go:build prod

package config

import "os"

func Load() *Config {
	cfg := loadBase()
	cfg.Environment = Prod
	cfg.DatabasePath = getEnv(CONFIG_DB_PATH, "")
	cfg.Username = getEnv(CONFIG_USERNAME, "")
	cfg.Password = getEnv(CONFIG_PASSWORD, "")

	os.Setenv(ENVIRONMENT, string(cfg.Environment))
	return cfg
}
