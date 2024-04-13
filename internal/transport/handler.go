package transport

import (
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/danielcosme/curious-ape/internal/application"
)

type Handler struct {
	App                  *application.App
	Version              string
	SessionManager       *scs.SessionManager
	templateCache        map[string]*template.Template
	partialTemplateCache map[string]*template.Template
}

func NewHandler(app *application.App, version string, sm *scs.SessionManager) (*Handler, error) {
	h := &Handler{
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
	h.templateCache = tc
	h.partialTemplateCache = tpc

	return h, nil
}

func (h *Handler) IsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedCtxKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
