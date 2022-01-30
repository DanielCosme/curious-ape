package transport

import (
	"fmt"
	"net/http"

	"github.com/danielcosme/curious-ape/internal/core/application"
)

// Transport is the top level object of the RESTApi transport layer
type Transport struct {
	App    *application.App
	Server *http.Server
}

func (a *Transport) ListenAndServe() error {
	fmt.Printf("Http server listening on: %s\n", a.Server.Addr)
	a.Server.Handler = a.Routes()
	return a.Server.ListenAndServe()
}
