package entity

import "time"

type IntegrationProvider string

const (
	ProviderFitbit = "fitbit"
	ProviderGoogle = "google"
	ProviderToggl  = "toggl"
	ProviderSelf   = "me"
)

type Oauth2 struct {
	ID           int                 `db:"id"`
	Provider     IntegrationProvider `db:"provider"`
	AccessToken  string              `db:"access_token"`
	RefreshToken string              `db:"refresh_token"`
	TokenType    string              `db:"token_type"`
	Expiration   time.Time           `db:"expiration"`
	// toggl extra
	ToogglWorkSpaceID    string `db:"toggl_workspace_id"`
	ToogglOrganizationID string `db:"toggl_organization_id"`
	ToogglProjectIDs     string `db:"toggl_project_ids"`
}

type Oauth2Config struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURL  string   `json:"redirect_url"`
	TokenURL     string   `json:"token_url"`
	AuthURL      string   `json:"auth_url"`
	Scopes       []string `json:"scopes"`
}

type Oauth2Filter struct {
	ID       []int
	Provider []IntegrationProvider
}
