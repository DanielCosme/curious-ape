package web

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func ChiRoutes(h *Handler) http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Handle("/static/", http.StripPrefix("/static/", fileServer))

	r.Get("/", h.home)
	r.Get("/habit/view", h.habitView)
	r.Post("/habit/create", h.habitCreate)

	return r
}
