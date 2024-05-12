package core

type Integration string

const (
	IntegrationFitbit = "fitbit"
	IntegrationGoogle = "google"
	IntegrationToggl  = "toggl"
	IntegrationSelf   = "me"
)

type Oauth2Config struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURL  string   `json:"redirect_url"`
	TokenURL     string   `json:"token_url"`
	AuthURL      string   `json:"auth_url"`
	Scopes       []string `json:"scopes"`
}
