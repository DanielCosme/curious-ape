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
		r.Use(a.authenticate)
		r.Get("/fitbit/authorize", a.authorizeFitbitHandler)
		r.Get("/fitbit/success", a.successFitbitHandler)
		r.Get("/google/authorize", a.authorizeGoogleHandler)
		r.Get("/google/success", a.successGoogleHandler)

		r.Get("/sleep/logs", a.listSleepRecordsHandler)
		r.Get("/sleep/log/{date}", a.showSleepRecordHandler)

		r.Get("/work/logs", a.listWorkRecordsHandler)
		r.Get("/work/log/{date}", a.showWorkRecordHandler)

		r.Get("/habit/{id}", a.showHabitHandler)
		r.Put("/habit/{id}", a.updateHabitHandler)
		r.Delete("/habit/{id}", a.deleteHabitHandler)
		r.Post("/habits", a.createHabitHandler)
		r.Get("/habits", a.listHabitsHandler)

		r.Post("/users", a.registerUserHandler)

		r.Get("/debug/seed", a.seedDataHandler)
		r.Get("/debug/misc", a.miscHandler)
		r.Get("/debug/day", a.dayGetHandler)

		r.Get("/healthcheck", a.healthcheckerHandler)
	})

	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}
