package routes

import (
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/danielcosme/curious-ape/internal/api/routes/middleware"
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/rest"
)

func ChiRoutes(a *application.App) http.Handler {
	h := Handler{App: a}
	r := chi.NewRouter()

	//  GET 	/habits 				-> all habits
	//  GET 	/habits/{id}			-> habit by ID
	//	GET 	/habits/date/{date} 	-> get all habits for day
	//  PUT 	/habits/{id} 			-> habit by ID
	//  DELETE 	/habits/{id} 			-> habit by ID
	//  POST 	/habits/date/{date} 	-> create habit for date

	r.Use(middleware.Logger(a))
	r.Use(middleware.RecoverPanic(a))
	r.Use(rest.MiddlewareParseForm)

	r.Get("/ping", h.Ping)

	// Days
	r.Route("/days", func(r chi.Router) {
		r.Get("/", h.DaysGetAll)
	})
	// Habits
	r.Route("/habits", func(r chi.Router) {
		r.Get("/", h.HabitsGetAll)
		r.Get("/categories", h.HabitsGetCategories)
		r.With(middleware.SetDay(a)).Post("/date/{date}", h.HabitCreate)
		r.Route("/{habitID}", func(r chi.Router) {
			r.Use(middleware.SetHabit(a))
			r.Get("/", h.HabitGet)
			r.Put("/", h.HabitUpdate)
			r.Delete("/", h.HabitDelete)
		})
	})
	// Sleep
	r.Route("/sleep", func(r chi.Router) {
		r.Get("/", h.SleepGetAll)
		r.With(middleware.SetDay(a)).Get("/date/{date}", h.SleepGetByDate)
	})
	// Sync
	r.Route("/sync", func(r chi.Router) {
		r.Post("/sleep/date/{start}", h.SyncSleepByDate)
		r.Post("/sleep/date/{startDate}/{endDate}", h.SyncSleepByDateRange)
	})
	// Oauth2
	r.Route("/oauth2/{provider}", func(r chi.Router) {
		r.Get("/connect", h.Oauth2Connect)
		r.Get("/success", h.Oauth2Success)
	})

	r.NotFound(h.NotFound)
	return r
}
