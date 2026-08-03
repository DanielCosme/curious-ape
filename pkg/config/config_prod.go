//go:build prod

package config

import "os"

func Load() *Config {
	cfg := loadBase()
	cfg.Environment = Prod

	os.Setenv(ENVIRONMENT, string(cfg.Environment))
	return cfg
}
