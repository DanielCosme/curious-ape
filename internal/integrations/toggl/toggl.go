package toggl

import (
	"io"
	"net/http"
)

type API struct {
	Me *MeService
}

func NewApi(token string, out io.Writer) *API {
	a := &API{
		Me: &MeService{client: &Client{Client: http.DefaultClient, token: token, out: out}},
	}
	return a
}
