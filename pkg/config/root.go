package config

import "strings"

// Configuration keys
const (
	CONFIG_PORT      = "APE_PORT"
	CONFIG_LOG_LEVEL = "APE_LOG_LEVEL" // DEBUG, INFO, WARN, ERROR
	CONFIG_DB_PATH   = "APE_DB_PATH"
	CONFIG_USERNAME  = "APE_USERNAME"
	CONFIG_PASSWORD  = "APE_PASSWORD"

	CONFIG_HEVY_API_KEY = "HEVY_API_KEY"

	CONFIG_TOGGL_API_KEY      = "TOGGL_API_KEY"
	CONFIG_TOGGL_WORKSPACE_ID = "TOGGL_WORKSPACE_ID"

	CONFIG_FITBIT_CLIENT_ID     = "FITBIT_CLIENT_ID"
	CONFIG_FITBIT_CLIENT_SECRET = "FITBIT_CLIENT_SECRET"
	CONFIG_FITBIT_REDIRECT_URL  = "FITBIT_REDIRECT_URL"
	CONFIG_FITBIT_TOKEN_URL     = "FITBIT_TOKEN_URL"
	CONFIG_FITBIT_AUTH_URL      = "FITBIT_AUTH_URL"
	CONFIG_FITBIT_SCOPES        = "FITBIT_SCOPES"
)

const (
	APP_NAME               = "ape"
	ENVIRONMENT            = "APE_ENVIRONMENT"
	MIGRATIONS_LOCATION    = "database/migrations/sqlite"
	DEPLOYMENT_DIR         = "deployment"
	PROD_USER              = "daniel"
	PROD_ADMIN             = "arch"
	REGISTRY               = "danicos.dev"
	TMP_DIR                = "./tmp"
	DATASTAR               = "https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.2/bundles/datastar.js"
	DOWNLOAD_DATASTAR      = "https://raw.githubusercontent.com/starfederation/datastar/develop/bundles/datastar.js"
	DOWNLOAD_DATASTAR_MAP  = "https://raw.githubusercontent.com/starfederation/datastar/develop/bundles/datastar.js.map"
	TZ                     = "America/Toronto"
	KUBERNETES_NAME        = "curious-ape"
	KUBERNETES_PORT        = 4000
	KUBERNETES_HOST        = "ape.danicos.me"
	KUBERNETES_DEPLOYMENT  = DEPLOYMENT_DIR + "/kubernetes"
	KUBERNETES_ENC_SECRETS = KUBERNETES_DEPLOYMENT + "/overlays/config"
	KUBERNETES_SECRETS     = TMP_DIR + "/secrets"
	LITESTREAM_IMAGE       = "docker.io/litestream/litestream:0.5.11"
)

var (
	KUBERNETES_IMAGE = "danicos.dev/daniel/curious-ape:"
)

func init() {
	KUBERNETES_IMAGE += strings.TrimPrefix(VERSION, "v")
}
