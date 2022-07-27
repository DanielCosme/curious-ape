package google

import (
	"net/http"
	"net/url"
	"time"
)

type FitnessService struct {
	client Client
}

type SessionsEnvelope struct {
	Session        []Session `json:"session,omitempty"`
	DeletedSession []Session `json:"deletedSession,omitempty"`
	NextPageToken  string    `json:"nextPageToken,omitempty"`
	HasMoreData    bool      `json:"hasMoreData,omitempty"`
}
type Application struct {
	PackageName string `json:"packageName,omitempty"`
	Version     string `json:"version,omitempty"`
	DetailsURL  string `json:"detailsUrl,omitempty"`
	Name        string `json:"name,omitempty"`
}
type Session struct {
	ID                 string      `json:"id,omitempty"`
	Name               string      `json:"name,omitempty"`
	Description        string      `json:"description,omitempty"`
	StartTimeMillis    string      `json:"startTimeMillis,omitempty"`
	EndTimeMillis      string      `json:"endTimeMillis,omitempty"`
	ModifiedTimeMillis string      `json:"modifiedTimeMillis,omitempty"`
	ActivityType       int         `json:"activityType,omitempty"`
	ActiveTimeMillis   string      `json:"activeTimeMillis,omitempty"`
	Application        Application `json:"application,omitempty"`
}

func (s *FitnessService) GetFitnessSessions(startTime, endTime time.Time) ([]Session, error) {
	var se *SessionsEnvelope
	params := url.Values{}
	params.Set("startTime", formatTime(startTime))
	params.Set("endTime", formatTime(endTime))
	err := s.client.Call(http.MethodGet, "/fitness/v1/users/me/sessions", params, &se)
	if err != nil {
		return nil, err
	}
	return se.Session, nil
}

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}
