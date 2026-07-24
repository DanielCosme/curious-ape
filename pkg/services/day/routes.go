package day

import (
	"danicos.dev/daniel/curious-ape/pkg/event"
	"github.com/go-chi/chi/v5"
	"github.com/stephenafamo/bob"
)

func SetupRoutes(r chi.Router, db bob.DB, bus event.Broker) error {
	svc := NewService(db, bus)
	handler := NewHandler(svc)

	r.Get("/", handler.Index)

	return nil
}
