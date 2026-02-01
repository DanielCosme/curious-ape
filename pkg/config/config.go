package config

const (
	APP     = "ape"
	VERSION = "v0.21.1"
)

const (
	ENVIRONMENT         = "APE_ENVIRONMENT"
	MIGRATIONS_LOCATION = "database/migrations/sqlite"
	DEPLOYMENT_DIR      = "deployment"
	PROD_HOST           = "ape-0" // Tailscale hostname (SSH in the server only works with VPN)
	PROD_USER           = "daniel"
	PROD_ADMIN          = "arch"
)
