package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/danielcosme/curious-ape/internal/api/routes/middleware"
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/rest"
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
	// r.Use(rest.AllowCORS) TODO allow when trusted origins is implemented

	r.Route("/api", func(r chi.Router) {
		r.Get("/ping", h.Ping)

		// Days
		r.Route("/days", func(r chi.Router) {
			r.Get("/", h.DaysGetAll)
			r.Route("/date/{date}", func(r chi.Router) {
				r.Use(middleware.SetDay(a))
				r.Get("/", h.DayGetByDate)
				r.Put("/", h.DayUpdate)
			})
		})

		r.Get("/tester", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("<p>Hello Daniel<p>"))
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
			r.Route("/date/{date}", func(r chi.Router) {
				r.Use(middleware.SetDay(a))
				r.Get("/", h.HabitsGetByDay)
				r.Post("/", h.HabitCreate)
			})
		})
		// Sleep
		r.Route("/sleep", func(r chi.Router) {
			r.Get("/", h.SleepGetAll)
			r.Route("/date/{date}", func(r chi.Router) {
				r.Use(middleware.SetDay(a))
				r.Get("/", h.SleepGetForDate)
				r.Post("/", h.SleepCreate)
			})
			r.Route("/{id}", func(r chi.Router) {
				r.Use(middleware.SetSleepLog(a))
				r.Get("/", h.SleepGet)
				r.Put("/", h.SleepUpdate)
				r.Delete("/", h.SleepDelete)
			})
		})
		// Fitness
		r.Route("/fitness", func(r chi.Router) {
			r.Get("/", h.FitnessGetAll)
			r.Route("/date/{date}", func(r chi.Router) {
				r.Use(middleware.SetDay(a))
				r.Get("/", h.FitnessGetForDate)
				r.Post("/", h.FitnessCreate)
			})
			r.Route("/{id}", func(r chi.Router) {
				r.Use(middleware.SetFitnessLog(a))
				r.Get("/", h.FitnessGet)
				r.Put("/", h.FitnessUpdate)
				r.Delete("/", h.FitnessDelete)
			})
		})
		// sync
		r.Route("/sync/{resourceToSync}", func(r chi.Router) {
			r.Post("/", h.Sync)
			r.Post("/date/{start}", h.SyncByDate)
			r.Post("/date/{startDate}/{endDate}", h.SyncByDateRange)
		})
		// Oauth2
		r.Route("/oauth2/{provider}", func(r chi.Router) {
			r.Get("/connect", h.Oauth2Connect)
			r.Get("/success", h.Oauth2Success)
			r.Post("/addToken", h.AddToken)
		})
		// Provider specific endpoints
		r.Route("/providers/toggl", func(r chi.Router) {
			r.Get("/projects", h.TogglGetProjects)
			r.Post("/projects", h.TogglAssignProjectsToGoal)
		})

		r.NotFound(h.NotFound)
	})

	r.NotFound(h.NotFound)
	return r
}
