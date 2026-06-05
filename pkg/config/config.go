package config

import "strings"

const (
	APP_NAME              = "ape"
	ENVIRONMENT           = "APE_ENVIRONMENT"
	MIGRATIONS_LOCATION   = "database/migrations/sqlite"
	DEPLOYMENT_DIR        = "deployment"
	PROD_USER             = "daniel"
	PROD_ADMIN            = "arch"
	REGISTRY              = "danicos.dev"
	DATASTAR              = "https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.2/bundles/datastar.js"
	TZ                    = "America/Toronto"
	KUBERNETES_NAME       = "curious-ape"
	KUBERNETES_PORT       = 4000
	KUBERNETES_HOST       = "ape.danicos.me"
	KUBERNETES_DEPLOYMENT = DEPLOYMENT_DIR + "/kubernetes"
	LITESTREAM_IMAGE      = "docker.io/litestream/litestream:0.5.11-scratch"
)

var (
	KUBERNETES_IMAGE = "danicos.dev/daniel/curious-ape:"
)

func init() {
	KUBERNETES_IMAGE += strings.TrimPrefix(VERSION, "v")
}
