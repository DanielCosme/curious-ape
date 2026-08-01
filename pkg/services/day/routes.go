package day

import (
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r chi.Router, handler *Handler) error {
	r.Get("/", handler.index)
	r.Route("/days", func(r chi.Router) {
		r.Get("/stream", handler.streamSSE)
		r.Post("/{date}/sync", handler.sync)
	})

	return nil
}
