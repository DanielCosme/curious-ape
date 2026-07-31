package habit

import (
	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
	"github.com/stephenafamo/bob"
)

func SetupRoutes(r chi.Router, db bob.DB, ns *nats.Conn) error {
	svc := NewService(db, ns)
	handler := NewHandler(svc)

	r.Route("/habits", func(r chi.Router) {
		r.Get("/", handler.Habits)
		r.Put("/{id}/flip", handler.Flip)
	})

	return nil
}
