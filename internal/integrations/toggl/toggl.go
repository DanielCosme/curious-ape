package toggl

import (
	"net/http"
)

type API struct {
	Me        *MeService
	Reports   *ReportsService
	Workspace *WorkspaceService
	Projects  *ProjectsService
}

func NewApi(workspaceID int, token string) *API {
	client := &Client{
		Client:      http.DefaultClient,
		workspaceID: workspaceID,
		token:       token,
	}
	a := &API{
		Me:        &MeService{client: client},
		Reports:   &ReportsService{client: client},
		Workspace: &WorkspaceService{client: client},
		Projects:  &ProjectsService{client: client},
	}
	return a
}
