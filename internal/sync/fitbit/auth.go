package fitbit

import (
	"os"

	"github.com/danielcosme/curious-ape/internal/auth"
)

var fitbitAuth = &auth.AuthConfig{}

func init() {
	fitbitAuth.AuthorizationURL = "https://www.fitbit.com/oauth2/authorize"
	fitbitAuth.TokenRequestURL = "https://api.fitbit.com/oauth2/token"
	fitbitAuth.ClientSecret = os.Getenv("FITBIT_SECRET")
	fitbitAuth.ClientID = os.Getenv("FITBIT_ID")
	fitbitAuth.RedirectURL = os.Getenv("FITBIT_REDIRECT_URI")
}
