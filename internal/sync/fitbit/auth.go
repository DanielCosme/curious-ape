package fitbit

import (
	"os"

	"github.com/danielcosme/curious-ape/internal/auth"
)

var FitbitAuth = &auth.AuthConfig{}

func init() {
	FitbitAuth.AuthorizationURL = "https://www.fitbit.com/oauth2/authorize"
	FitbitAuth.TokenRequestURL = "https://api.fitbit.com/oauth2/token"
	FitbitAuth.ClientSecret = os.Getenv("FITBIT_SECRET")
	FitbitAuth.ClientID = os.Getenv("FITBIT_ID")
	FitbitAuth.RedirectURL = os.Getenv("FITBIT_REDIRECT_URI")
	FitbitAuth.Provider = "fitbit"
}
