package habit

import (
	"github.com/go-chi/chi/v5"
	"github.com/stephenafamo/bob"
)

func SetupRoutes(r chi.Router, db bob.DB) error {
	svc := NewService(db)
	handler := NewHandler(svc)

	r.Route("/habits", func(r chi.Router) {
		r.Get("/", handler.Habits)
	})

	return nil
}
