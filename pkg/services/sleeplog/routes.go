package sleeplog

import (
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r chi.Router, handler *Handler) error {
	r.Route("/sleep", func(r chi.Router) {
		r.Get("/", handler.sleeplogPage)
	})
	return nil
}
