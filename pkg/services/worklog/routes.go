package worklog

import (
	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
	"github.com/stephenafamo/bob"
)

func SetupRoutes(r chi.Router, db bob.DB, ns *nats.Conn) error {
	svc := NewService(db, ns)
	handler := NewHandler(svc)

	r.Route("/worklog", func(r chi.Router) {
		r.Post("/", handler.worklogPage)
	})

	return nil
}
