package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/danielcosme/curious-ape/internal/core/application"
)

type Handler struct {
	App           *application.App
	templateCache map[string]*template.Template
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
