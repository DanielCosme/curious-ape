package toggl

import (
	"fmt"
	"net/http"
	"time"
)

type ReportsService struct {
	client *Client
}

type SummaryEnvelope struct {
	Logs          []*Summary
	TotalDuration time.Duration
}

type Summary struct {
	BillableSeconds int `json:"billable_seconds"`
	ProjectID       int `json:"project_id"`
	TrackedSeconds  int `json:"tracked_seconds"`
	UserID          int `json:"user_id"`
}

func (s *ReportsService) GetDaySummary(day time.Time) (res SummaryEnvelope, err error) {
	type query struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	d := FormatDate(day)
	ss, err := s.summaryRequest(&query{
		StartDate: d,
		EndDate:   d,
	})
	if err != nil {
		return
	}
	totalSeconds := 0
	for _, s := range ss {
		totalSeconds += s.TrackedSeconds
	}
	res.TotalDuration = time.Second * time.Duration(totalSeconds)
	return
}

func (s *ReportsService) summaryRequest(body any) ([]*Summary, error) {
	var summary []*Summary
	err := s.client.Call(http.MethodPost, fmt.Sprintf("/reports/api/v3/workspace/%d/projects/summary", s.client.workspaceID), body, &summary)
	return summary, err
}

func FormatDate(time time.Time) string {
	return time.Format("2006-01-02")
}
