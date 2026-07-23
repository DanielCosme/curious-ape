package api

import (
	"net/http"

	"danicos.dev/daniel/curious-ape/pkg/application"
	"github.com/alexedwards/scs/v2"
)

type ContextKey string

const (
	ctxKeyIsAuthenticated     ContextKey = "isAuthenticated"
	ctxKeyAuthenticatedUserID ContextKey = "authenticatedUserID"
	ctxUser                   ContextKey = "user"
)

type OldAPI struct {
	App     *application.App
	Scs     *scs.SessionManager
	Version string
}

func NewApi(app *application.App, sessionManager *scs.SessionManager, version string) *OldAPI {
	return &OldAPI{
		App:     app,
		Version: version,
		Scs:     sessionManager,
	}
}

func (a *OldAPI) IsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(ctxKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
