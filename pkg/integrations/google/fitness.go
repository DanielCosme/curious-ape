package google

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/danielcosme/curious-ape/pkg/core"
)

type FitnessService struct {
	client Client
}

type SessionsEnvelope struct {
	Session        []FitnessSession `json:"session,omitempty"`
	DeletedSession []FitnessSession `json:"deletedSession,omitempty"`
	NextPageToken  string           `json:"nextPageToken,omitempty"`
	HasMoreData    bool             `json:"hasMoreData,omitempty"`
}
type Application struct {
	PackageName string `json:"packageName,omitempty"`
	Version     string `json:"version,omitempty"`
	DetailsURL  string `json:"detailsUrl,omitempty"`
	Name        string `json:"name,omitempty"`
}
type FitnessSession struct {
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

func (s *FitnessService) GetFitnessSessions(startTime, endTime time.Time) ([]FitnessSession, error) {
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

func ParseMillis(s string) time.Time {
	millis, _ := strconv.Atoi(s)
	t := time.UnixMilli(int64(millis))
	return core.TimeUTC(t)
}
