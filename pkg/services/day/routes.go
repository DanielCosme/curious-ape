package day

import (
	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
	"github.com/stephenafamo/bob"
)

func SetupRoutes(r chi.Router, db bob.DB, ns *nats.Conn) error {
	svc := NewService(db, ns)
	handler := NewHandler(svc)

	r.Get("/", handler.index)
	r.Route("/days", func(r chi.Router) {
		r.Get("/stream", handler.streamSSE)
		r.Post("/{date}/sync", handler.sync)
	})

	return nil
}
