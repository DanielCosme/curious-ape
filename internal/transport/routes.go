package transport

import (
	"net/http"

	"github.com/danielcosme/curious-ape/web"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func ChiRoutes(h *Handler) http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(h.midRecoverPanic)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(midSecureHeaders)

	fileServer := http.FileServer(http.FS(web.Files))
	r.Handle("/static/*", fileServer)

	r.Route("/login", func(r chi.Router) {
		r.Use(h.SessionManager.LoadAndSave)

		r.Get("/", h.loginForm)
		r.Post("/", h.loginPost)
	})

	r.Get("/api/oauth2/{provider}/success", h.Oauth2Success)
	// Protected routes.
	r.Route("/", func(r chi.Router) {
		r.Use(h.SessionManager.LoadAndSave)
		r.Use(h.midAuthenticate)
		r.Use(h.RequireAuth)

		r.Get("/", h.home)
		r.Post("/logout", h.logout)

		// Habits.
		r.With(h.midSetHabit).Get("/habit/view/{id}", h.habit)
		r.Get("/habit/new", h.newHabitForm)
		r.Post("/habit/new", h.newHabitPost)
		r.Post("/habit/log", h.newHabitLogPost)

		// Oauth2
		r.Route("/oauth2/{provider}", func(r chi.Router) {
			r.Get("/connect/form", h.Oauth2ConnectForm)
			r.Get("/connect", h.Oauth2Connect)
			r.Post("/addToken", h.AddToken)
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) { h.notFound(w) })
	return r
}
