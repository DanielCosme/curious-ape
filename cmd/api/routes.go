package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.NotFound(a.notFoundResponse)
	mux.MethodNotAllowed(a.methodNotAllowedResponse)

	mux.Use(a.recoverPanic)

	mux.Route("/v1", func(r chi.Router) {
		r.Get("/healthcheck", a.healthcheckerHandler)

		r.Get("/food/habit/{date}", a.showFoodHabitHandler)
		r.Put("/food/habit/{date}", a.updateFoodHabitHandler)
		r.Delete("/food/habit/{id}", a.deleteFoodHabitHandler)
		r.Post("/food/habits", a.createFoodHabitHandler)
		r.Get("/food/habits", a.listFoodHabitsHandler)
	})

	return mux
}
