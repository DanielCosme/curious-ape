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

func NewApi(token string) *API {
	a := &API{
		Me:        &MeService{client: &Client{Client: http.DefaultClient}},
		Reports:   &ReportsService{client: &Client{Client: http.DefaultClient}},
		Workspace: &WorkspaceService{client: &Client{Client: http.DefaultClient}},
		Projects:  &ProjectsService{client: &Client{Client: http.DefaultClient}},
	}
	return a
}
