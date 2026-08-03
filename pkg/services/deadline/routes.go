package deadline

import (
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r chi.Router, handler *Handler) error {
	r.Route("/deadlines", func(r chi.Router) {
		r.Get("/", handler.deadlinePage)
		r.Get("/new", handler.newDeadlinePage)
		r.Post("/new", handler.newDeadlinePost)
	})
	return nil
}
