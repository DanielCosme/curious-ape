package web

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func ChiRoutes(h *Handler) http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(h.midRecoverPanic)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(midSecureHeaders)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	r.Get("/", h.home)

	// Habits.
	r.Route("/habit", func(r chi.Router) {
		r.With(h.midSetHabit).Get("/view/{id}", h.habitView)
		r.Get("/create", h.habitCreateForm)
		r.Post("/create", h.habitCreate)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) { h.notFound(w) })
	return r
}
