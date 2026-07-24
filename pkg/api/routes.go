package api

import (
	"context"
	"io"
	"net/http"

	"github.com/stephenafamo/bob"

	"danicos.dev/daniel/curious-ape/pkg/config"
	"danicos.dev/daniel/curious-ape/pkg/services/user"
	"danicos.dev/daniel/curious-ape/pkg/web"
	"danicos.dev/daniel/curious-ape/web/resources"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(ctx context.Context, r chi.Router, sessionManager *scs.SessionManager, db bob.DB) error {
	if config.Global.Environment == config.Dev {
		setupReload(r)
	}

	r.Handle("/static/*", resources.Handler())

	r.Group(func(r chi.Router) {
		r.Use(sessionManager.LoadAndSave)
		r.Use(user.MiddlewareAuthenticateFromSession(sessionManager, db))

		user.SetupRoutes(r, db, sessionManager)
	})

	return nil
}

func Index(w http.ResponseWriter, r *http.Request) {
	// io.WriteString(w, "Index")
	web.Redirect(w, "/login")
}

func putHandler(w http.ResponseWriter, r *http.Request, sessionManager *scs.SessionManager) {
	// Store a new key and value in the session data.
	sessionManager.Put(r.Context(), "message", "Hello from a session!")
}

func getHandler(w http.ResponseWriter, r *http.Request, sessionManager *scs.SessionManager) {
	// Use the GetString helper to retrieve the string value associated with a
	// key. The zero value is returned if the key does not exist.
	msg := sessionManager.GetString(r.Context(), "message")
	io.WriteString(w, msg)
}

func setupReload(router chi.Router) {
	// TODO: Implement Hot Reload
}

func mid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
