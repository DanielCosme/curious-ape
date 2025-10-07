package toggl

import (
	"github.com/danielcosme/curious-ape/pkg/core"
	"net/http"
	"net/url"
	"time"
)

type TimeEntries struct {
	client *Client
	tz     *time.Location
}

type TimeEntry struct {
	ID              int64     `json:"id"`
	WorkspaceID     int       `json:"workspace_id"`
	ProjectID       int       `json:"project_id"`
	TaskID          any       `json:"task_id"`
	Billable        bool      `json:"billable"`
	Start           time.Time `json:"start"`
	Stop            time.Time `json:"stop"`
	Duration        int       `json:"duration"`
	Description     string    `json:"description"`
	Tags            []string  `json:"tags"`
	TagIds          []int     `json:"tag_ids"`
	Duronly         bool      `json:"duronly"`
	At              time.Time `json:"at"`
	ServerDeletedAt any       `json:"server_deleted_at"`
	UserID          int       `json:"user_id"`
	UID             int       `json:"uid"`
	Wid             int       `json:"wid"`
	Pid             int       `json:"pid"`
}

func (te *TimeEntries) GetDayEntries(date time.Time) (tes []TimeEntry, err error) {
	query := url.Values{}
	query.Add("start_date", date.Format("2006-01-02"))
	query.Add("end_date", date.AddDate(0, 0, 2).Format("2006-01-02"))
	path := "/api/v9/me/time_entries?" + query.Encode()
	err = te.client.Call(http.MethodGet, path, nil, &tes)
	if err == nil {
		for idx, timeEntry := range tes {
			tes[idx].Start = NormalizeLocation(timeEntry.Start, te.tz)
			tes[idx].Stop = NormalizeLocation(timeEntry.Stop, te.tz)
			tes[idx].At = NormalizeLocation(timeEntry.At, te.tz)
		}
	}
	return
}

func NormalizeLocation(t time.Time, loc *time.Location) time.Time {
	return core.TimeUTC(t.In(loc))
}
