package main

import (
	"expvar"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.NotFound(a.notFoundResponse)
	mux.MethodNotAllowed(a.methodNotAllowedResponse)

	mux.Use(a.metrics)
	mux.Use(a.recoverPanic)
	mux.Use(a.rateLimit)

	mux.Route("/v1", func(r chi.Router) {
		r.Get("/healthcheck", a.healthcheckerHandler)

		r.Get("/food/habit/{date}", a.showFoodHabitHandler)
		r.Put("/food/habit/{date}", a.updateFoodHabitHandler)
		r.Delete("/food/habit/{id}", a.deleteFoodHabitHandler)
		r.Post("/food/habits", a.createFoodHabitHandler)
		r.Get("/food/habits", a.listFoodHabitsHandler)

		r.Post("/users", a.registerUserHandler)
	})

	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}
