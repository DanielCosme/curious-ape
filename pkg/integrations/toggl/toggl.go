package toggl

import (
	"net/http"
	"time"
)

type API struct {
	Me          *MeService
	Reports     *ReportsService
	Workspace   *WorkspaceService
	Projects    *ProjectsService
	TimeEntries *TimeEntries
	timeZone    *time.Location
}

func NewApi(workspaceID int, token string) (*API, error) {
	client := &Client{
		Client:      http.DefaultClient,
		workspaceID: workspaceID,
		token:       token,
	}
	a := &API{
		Me:          &MeService{client: client},
		Reports:     &ReportsService{client: client},
		Workspace:   &WorkspaceService{client: client},
		Projects:    &ProjectsService{client: client},
		TimeEntries: &TimeEntries{client: client},
	}

	profile, err := a.Me.GetProfile()
	if err != nil {
		return nil, err
	}
	err = a.SetClientTimeZone(profile.Timezone)
	if err != nil {
		return nil, err
	}
	a.TimeEntries.tz = a.ClientTimezone()
	return a, nil
}

func (a *API) SetClientTimeZone(tz string) error {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return err
	}
	a.timeZone = loc
	a.TimeEntries.tz = loc

	return nil
}

func (a *API) ClientTimezone() *time.Location {
	return a.timeZone
}
