package habit

import (
	"danicos.dev/daniel/curious-ape/pkg/event"
	"github.com/go-chi/chi/v5"
	"github.com/stephenafamo/bob"
)

func SetupRoutes(r chi.Router, db bob.DB, bus event.Broker) error {
	svc := NewService(db, bus)
	handler := NewHandler(svc)

	r.Route("/habits", func(r chi.Router) {
		r.Get("/", handler.Habits)
		r.Put("/{id}/flip", handler.Flip)
	})

	return nil
}
