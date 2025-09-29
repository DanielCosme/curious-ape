package api

import (
	"net/http"

	"github.com/danielcosme/curious-ape/pkg/tucan"
)

func EchoRoutes(a *API) http.Handler {
	t := tucan.New()

	// e.StaticFS("/static", echo.MustSubFS(views.StaticFS, "static"))

	return t
}
