//go:build dev || !prod

package config

func Load() *Config {
	cfg := loadBase()
	cfg.Environment = Dev
	return cfg
}
