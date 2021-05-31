package toggl

import (
	"os"

	"github.com/danielcosme/curious-ape/internal/auth"
)

var TogglAuth = &auth.AuthConfig{}

func init() {
	TogglAuth.AuthorizationURL = ""
	TogglAuth.TokenRequestURL = ""
	TogglAuth.ClientSecret = os.Getenv("TOGGL_PASS")
	TogglAuth.ClientID = os.Getenv("TOGGL_TOKEN")
	TogglAuth.RedirectURL = ""
	TogglAuth.Provider = "toggl"
}
