package fitbit

import (
	"os"

	"github.com/danielcosme/curious-ape/internal/auth"
)

type SleepCollector struct {
	Auth *auth.AuthConfig
	Name string
}

var fitbitAuth = &auth.AuthConfig{}
var Fitbit = &SleepCollector{
	Auth: fitbitAuth,
}

// TODO make fields inmutable by restricting write operations via method accessor
func init() {
	fitbitAuth.AuthorizationURL = "https://www.fitbit.com/oauth2/authorize"
	fitbitAuth.TokenRequestURL = "https://api.fitbit.com/oauth2/token"
	fitbitAuth.ClientSecret = os.Getenv("FITBIT_SECRET")
	fitbitAuth.ClientID = os.Getenv("FITBIT_ID")
	fitbitAuth.RedirectURL = os.Getenv("FITBIT_REDIRECT_URI")
}

// to load: secret, id
// redirect can redirect URI be hardcoded?? no, one for devel, one for prod.
