package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/danielcosme/curious-ape/internal/api/router/middleware"
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/rest"
)

func Routes(a *application.App) http.Handler {
	h := Handler{App: a}

	mux := http.NewServeMux()
	md := rest.NewMiddleware()

	md.Use(middleware.Logger(a))
	md.Use(middleware.RecoverPanic(a))
	md.Use(rest.MiddlewareParseForm)

	mux.Handle("/sleep", rest.NewMiddleware(middleware.SetDay(a)).Then(h.SleepLogs))
	mux.HandleFunc("/sleep/sync", h.SleepLogs)

	mux.HandleFunc("/oauth2/fitbit/connect", h.Oauth2FitbitConnect)
	mux.HandleFunc("/oauth2/fitbit/success", h.Oauth2FitbitSuccess)


	return md.Commit(mux)
}

func ChiRoutes(a *application.App) http.Handler {
	h := Handler{App: a}

	r := chi.NewRouter()

	//  POST /days/____/habits/create
	//	GET /days/___/habits 	-> all habits for day?

	//  GET /habits 		-> all habits
	//  GET /habits/1 		-> habit by ID
	//  PUT /habits/1 		-> habit by ID
	//  DELETE /habits/1 	-> habit by ID

	r.Use(middleware.Logger(a))
	r.Use(middleware.RecoverPanic(a))
	r.Use(rest.MiddlewareParseForm)

	r.Get("/ping", h.Ping)

	r.Route("/days", func(r chi.Router) {
		r.Get("/", h.DaysGetAll)
		r.Route("/{date}", func(r chi.Router) {
			r.Use(middleware.SetDay(a))
			r.Post("/habits", h.HabitCreate)
		})
	})
	r.Route("/habits", func(r chi.Router) {
		r.Get("/", h.HabitsGetAll)
		r.Get("/categories", h.HabitsGetCategories)
		r.Route("/{habitID}", func(r chi.Router) {
			r.Use(middleware.SetHabit(a))
			r.Get("/", h.HabitGet)
			r.Put("/", h.HabitUpdate)
			r.Delete("/", h.HabitDelete)
		})
	})

	r.NotFound(h.NotFound)
	return r
}
