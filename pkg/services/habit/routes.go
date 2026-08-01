package habit

import (
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r chi.Router, handler *Handler) error {
	r.Route("/habits", func(r chi.Router) {
		r.Get("/", handler.HabitsPage)
		r.Put("/{id}/flip", handler.Flip)
	})

	return nil
}
