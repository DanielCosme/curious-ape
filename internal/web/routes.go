package web

import (
	"fmt"
	"net/http"

	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	App *application.App
}

func ChiRoutes(a *application.App) http.Handler {
	h := Handler{App: a}
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<p>Hello, World</p>")
	})

	r.Get("/habit/view", h.habitView)
	r.Post("/habit/create", h.habitCreate)

	return r
}
