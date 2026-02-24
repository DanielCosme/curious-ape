package api

import (
	"net/http"

	"git.danicos.dev/daniel/curious-ape/pkg/application"
	"github.com/alexedwards/scs/v2"
)

type ContextKey string

const (
	ctxKeyIsAuthenticated     ContextKey = "isAuthenticated"
	ctxKeyAuthenticatedUserID ContextKey = "authenticatedUserID"
	ctxUser                   ContextKey = "user"
)

type API struct {
	App     *application.App
	Scs     *scs.SessionManager
	Version string
}

func NewApi(app *application.App, sessionManager *scs.SessionManager, version string) *API {
	return &API{
		App:     app,
		Version: version,
		Scs:     sessionManager,
	}
}

func (a *API) IsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(ctxKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
