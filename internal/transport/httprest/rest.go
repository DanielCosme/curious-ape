package httprest

import (
	"fmt"
	"github.com/danielcosme/curious-ape/rest"
	"net/http"

	"github.com/danielcosme/curious-ape/internal/core/application"
)

type API struct {
	App *application.App
	*http.Server
}

func (a *API) Run() error {
	fmt.Printf("Http server listening on: %s", a.Server.Addr)
	a.Handler = a.Routes()
	return a.ListenAndServe()
}

func (a *API) Ping(rw http.ResponseWriter, r *http.Request) {
	rest.JSON(rw, http.StatusOK, rest.Envelope{"message": "pong"})
}