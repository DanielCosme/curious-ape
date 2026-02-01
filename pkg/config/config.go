package config

const (
	APP_NAME            = "ape"
	ENVIRONMENT         = "APE_ENVIRONMENT"
	MIGRATIONS_LOCATION = "database/migrations/sqlite"
	DEPLOYMENT_DIR      = "deployment"
	PROD_HOST           = "ape-0" // Tailscale hostname (SSH in the server only works with VPN)
	PROD_USER           = "daniel"
	PROD_ADMIN          = "arch"
)
