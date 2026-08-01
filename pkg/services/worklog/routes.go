package worklog

import (
	"danicos.dev/daniel/curious-ape/pkg/core"
	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
	"github.com/stephenafamo/bob"
)

func SetupRoutes(r chi.Router, db bob.DB, ns *nats.Conn, integration core.WorkLogIntegration) error {
	svc := NewService(db, ns, integration)
	handler := NewHandler(svc)

	r.Route("/worklog", func(r chi.Router) {
		r.Get("/", handler.worklogPage)
	})

	return nil
}
