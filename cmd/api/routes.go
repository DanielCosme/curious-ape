package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Route("/v1", func(r chi.Router) {
		r.Get("/healthcheck", a.healthcheckerHandler)

		r.Get("/food/habit/{date}", a.showFoodHabitHandler)
		r.Post("/food/habits", a.createFoodHabitHandler)
	})

	mux.NotFound(a.notFoundResponse)
	mux.MethodNotAllowed(a.methodNotAllowedResponse)

	return mux
}
