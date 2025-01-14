package transport

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/danielcosme/curious-ape/internal/application"
)

type Transport struct {
	App            *application.App
	Version        string
	SessionManager *scs.SessionManager
}

func NewTransport(app *application.App, sm *scs.SessionManager, version string) *Transport {
	return &Transport{
		App:            app,
		Version:        version,
		SessionManager: sm,
	}
}

func (t *Transport) IsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(ctxKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
