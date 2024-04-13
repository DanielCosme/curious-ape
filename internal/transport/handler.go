package transport

import (
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/danielcosme/curious-ape/internal/application"
)

type Transport struct {
	App                  *application.App
	Version              string
	SessionManager       *scs.SessionManager
	templateCache        map[string]*template.Template
	partialTemplateCache map[string]*template.Template
}

func NewTransport(app *application.App, version string, sm *scs.SessionManager) (*Transport, error) {
	t := &Transport{
		App:            app,
		Version:        version,
		SessionManager: sm,
	}
	tc, err := newTemplateCache()
	if err != nil {
		return nil, err
	}
	tpc, err := newTemplatePartialCache()
	if err != nil {
		return nil, err
	}
	t.templateCache = tc
	t.partialTemplateCache = tpc

	return t, nil
}

func (h *Transport) IsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedCtxKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
