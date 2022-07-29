package toggl

import (
	"net/http"
	"time"
)

type WorkspaceService struct {
	client *Client
}

type Workspace struct {
	ID                          int         `json:"id"`
	OrganizationID              int         `json:"organization_id"`
	Name                        string      `json:"name"`
	Profile                     int         `json:"profile"`
	Premium                     bool        `json:"premium"`
	BusinessWs                  bool        `json:"business_ws"`
	Admin                       bool        `json:"admin"`
	SuspendedAt                 interface{} `json:"suspended_at"`
	ServerDeletedAt             interface{} `json:"server_deleted_at"`
	DefaultHourlyRate           interface{} `json:"default_hourly_rate"`
	RateLastUpdated             interface{} `json:"rate_last_updated"`
	DefaultCurrency             string      `json:"default_currency"`
	OnlyAdminsMayCreateProjects bool        `json:"only_admins_may_create_projects"`
	OnlyAdminsMayCreateTags     bool        `json:"only_admins_may_create_tags"`
	OnlyAdminsSeeBillableRates  bool        `json:"only_admins_see_billable_rates"`
	OnlyAdminsSeeTeamDashboard  bool        `json:"only_admins_see_team_dashboard"`
	ProjectsBillableByDefault   bool        `json:"projects_billable_by_default"`
	ReportsCollapse             bool        `json:"reports_collapse"`
	Rounding                    int         `json:"rounding"`
	RoundingMinutes             int         `json:"rounding_minutes"`
	APIToken                    string      `json:"api_token"`
	At                          time.Time   `json:"at"`
	LogoURL                     string      `json:"logo_url"`
	IcalURL                     string      `json:"ical_url"`
	IcalEnabled                 bool        `json:"ical_enabled"`
	CsvUpload                   interface{} `json:"csv_upload"`
	Subscription                interface{} `json:"subscription"`
}

func (s *WorkspaceService) Get() ([]*Workspace, error) {
	var workspaces []*Workspace
	err := s.client.Call(http.MethodGet, "/api/v9/workspaces", nil, &workspaces)
	return workspaces, err
}
