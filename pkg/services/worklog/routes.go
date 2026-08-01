package worklog

import (
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r chi.Router, handler *Handler) error {
	r.Route("/worklog", func(r chi.Router) {
		r.Get("/", handler.worklogPage)
	})
	return nil
}
