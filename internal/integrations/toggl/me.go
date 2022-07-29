package toggl

import (
	"net/http"
	"time"
)

type MeService struct {
	client *Client
}

type Me struct {
	ID                 int       `json:"id"`
	APIToken           string    `json:"api_token"`
	Email              string    `json:"email"`
	Fullname           string    `json:"fullname"`
	Timezone           string    `json:"timezone"`
	DefaultWorkspaceID int       `json:"default_workspace_id"`
	BeginningOfWeek    int       `json:"beginning_of_week"`
	ImageURL           string    `json:"image_url"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	OpenidEmail        string    `json:"openid_email"`
	OpenidEnabled      bool      `json:"openid_enabled"`
	CountryID          int       `json:"country_id"`
	At                 time.Time `json:"at"`
	IntercomHash       string    `json:"intercom_hash"`
	OauthProviders     []string  `json:"oauth_providers"`
	HasPassword        bool      `json:"has_password"`
}

func (s *MeService) GetProfile() (*Me, error) {
	var me *Me
	err := s.client.Call(http.MethodGet, "/api/v9/me", nil, &me)
	return me, err
}
