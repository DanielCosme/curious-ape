package transport

import (
	"fmt"
	"net/http"

	"github.com/danielcosme/curious-ape/internal/core/application"
)

// API is the top level object of the transport layer
type API struct {
	App    *application.App
	Server *http.Server
}

func (a *API) Run() error {
	fmt.Printf("Http server listening on: %s\n", a.Server.Addr)
	a.Server.Handler = a.Routes()
	return a.Server.ListenAndServe()
}
