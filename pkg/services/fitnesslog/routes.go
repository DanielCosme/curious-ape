package fitnesslog

import (
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r chi.Router, handler *Handler) error {
	r.Route("/fitness", func(r chi.Router) {
		r.Get("/", handler.fitnesslogPage)
	})
	return nil
}
