package toggl

import (
	"net/http"
	"net/url"
	"time"
)

type ReportsService struct {
	client *Client
}

type Summary struct {
	TotalGrand int    `json:"total_grand"`
	Data       []Data `json:"data"`
}
type Data struct {
	ID    int     `json:"id"`
	Title Title   `json:"title"`
	Time  int     `json:"time"`
	Items []Items `json:"items"`
}
type Title struct {
	Project  string `json:"project"`
	Color    string `json:"color"`
	HexColor string `json:"hex_color"`
}
type Items struct {
	Title      TimeEntryTitle `json:"title"`
	Time       int            `json:"time"`
	LocalStart string         `json:"local_start"`
}
type TimeEntryTitle struct {
	TimeEntry string `json:"time_entry"`
}

func (s *ReportsService) GetDaySummaryForProjectIDs(day time.Time, projectIDs, workspaceID string) (*Summary, error) {
	params := NewPrams(s.client.token, workspaceID)
	params.Set("since", FormatDate(day))
	params.Set("until", FormatDate(day))
	params.Set("project_ids", projectIDs)
	return s.summaryRequest(params)
}

// NewPrams returns mandatory params for the toggl api
func NewPrams(userAgent, workspaceID string) url.Values {
	params := url.Values{}
	params.Set("user_agent", userAgent)
	params.Set("workspace_id", workspaceID)
	return params
}

func (s *ReportsService) summaryRequest(params url.Values) (*Summary, error) {
	var summary *Summary
	err := s.client.Call(http.MethodGet, "/reports/api/v2/summary", params, &summary)
	return summary, err
}

func FormatDate(time time.Time) string {
	return time.Format("2006-01-02")
}
