package transport

import (
	"github.com/danielcosme/curious-ape/sdk/log"
	"net/http"

	"github.com/danielcosme/curious-ape/internal/core/application"
)

// API is the top level object of the transport layer
type API struct {
	App    *application.App
	Server *http.Server
}

func (api *API) Run() error {
	api.Server.Handler = api.Routes()

	api.App.Log.InfoP("HTTP server listening", log.Prop{"addr": api.Server.Addr})
	return api.Server.ListenAndServe()
}
