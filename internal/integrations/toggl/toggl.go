package toggl

import (
	"io"
	"net/http"
)

type API struct {
	Me        *MeService
	Reports   *ReportsService
	Workspace *WorkspaceService
}

func NewApi(token string, out io.Writer) *API {
	a := &API{
		Me:        &MeService{client: &Client{Client: http.DefaultClient, token: token, out: out}},
		Reports:   &ReportsService{client: &Client{Client: http.DefaultClient, token: token, out: out}},
		Workspace: &WorkspaceService{client: &Client{Client: http.DefaultClient, token: token, out: out}},
	}
	return a
}
