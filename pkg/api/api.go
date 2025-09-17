package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/danielcosme/curious-ape/pkg/application"
	"net/http"
)

type API struct {
	App            *application.App
	Version        string
	SessionManager *scs.SessionManager
}

func NewTransport(app *application.App, sm *scs.SessionManager, version string) *API {
	return &API{
		App:            app,
		Version:        version,
		SessionManager: sm,
	}
}

func (api *API) IsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(ctxKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
