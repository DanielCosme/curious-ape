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

	//  GET 	/habits 				-> get all habits
	//  GET 	/habits/{id}			-> get habit by ID
	//	GET 	/habits/date/{date} 	-> get all habits for day
	//  PUT 	/habits/{id} 			-> update habit by ID
	//  DELETE 	/habits/{id} 			-> delete habit by ID
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
		r.Route("/{habitID}", func(r chi.Router) {
			r.Use(middleware.SetHabit(a))
			r.Get("/", h.HabitGet)
			r.Put("/", h.HabitUpdate)
			r.Delete("/", h.HabitDelete)
		})
		r.With(middleware.SetDay(a)).Post("/date/{date}", h.HabitCreate)
	})
	// Sleep
	r.Route("/sleep", func(r chi.Router) {
		r.Get("/", h.SleepGetAll)
		r.Delete("/{id}", h.SleepDeleteByID)
		r.Route("/date/{date}", func(r chi.Router) {
			r.Use(middleware.SetDay(a))
			r.Get("/", h.SleepGetByDate)
			r.Post("/", h.SleepCreate)
		})
	})
	// Fitness
	r.Route("/fitness", func(r chi.Router) {
		r.Get("/", h.FitnessGet)
		r.Route("/date/{date}", func(r chi.Router) {
			r.Use(middleware.SetDay(a))
			r.Get("/", h.FitnessGetByDate)
			// r.Post("/", h.SleepCreate)
		})
	})
	// Sync
	r.Route("/sync", func(r chi.Router) {
		r.Post("/sleep", h.SyncSleep)
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
