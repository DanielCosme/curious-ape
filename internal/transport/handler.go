package transport

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/danielcosme/curious-ape/internal/core/application"
)

type Handler struct {
	App                  *application.App
	Version              string
	SessionManager       *scs.SessionManager
	templateCache        map[string]*template.Template
	partialTemplateCache map[string]*template.Template
}

func (h *Handler) IsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedCtxKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}
