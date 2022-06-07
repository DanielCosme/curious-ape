package api

import (
	"github.com/danielcosme/curious-ape/sdk/log"
	"net/http"

	"github.com/danielcosme/curious-ape/internal/core/application"
)

// Transport is the top level object of the API layer
type Transport struct {
	App    *application.App
	Server *http.Server
}

func (t *Transport) Run() error {
	t.Server.Handler = t.Routes()

	t.App.Log.InfoP("HTTP server listening", log.Prop{"addr": t.Server.Addr})
	return t.Server.ListenAndServe()
}
