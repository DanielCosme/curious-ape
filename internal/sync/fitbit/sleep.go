package fitbit

import (
	"github.com/danielcosme/curious-ape/internal/auth"
)

type SleepCollector struct {
	Auth *auth.AuthConfig
	name string
}

var Fitbit = &SleepCollector{
	Auth: fitbitAuth,
}

func (col *SleepCollector) AuthorizationURI() string {
	params := map[string]string{
		"client_id":     col.Auth.ClientID,
		"response_type": "code",
		"scope":         "sleep",
		"expires_in":    "604800",
		"redirect_uri":  col.Auth.RedirectURL,
	}
	urlEncoded := col.Auth.AuthorizationURL + "?" + auth.UrlEncode(params)
	return urlEncoded
}
