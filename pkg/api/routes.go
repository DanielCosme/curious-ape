package api

import (
	"context"
	"database/sql"
	"io"
	"net/http"

	"danicos.dev/daniel/curious-ape/pkg/config"
	"danicos.dev/daniel/curious-ape/web/resources"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(ctx context.Context, r chi.Router, sessionManager *scs.SessionManager, db *sql.DB) error {
	if config.Global.Environment == config.Dev {
		setupReload(r)
	}

	r.Handle("/static/*", resources.Handler())

	r.Group(func(r chi.Router) {
		r.Use(sessionManager.LoadAndSave)

		r.Get("/", Index)
		r.Get("/login", LoginPage)
	})

	// Implement Hotreload?

	// TODO: Return the Login Page.
	// ALSO: If Non-Authenticated, Redirect to /Login
	//
	// 2. Authenticate from session
	// 3. Set UI State?
	//

	// Here, I stablish the protected from the non-protected routes.
	// So, the first thing I need to do is to Return the Login Page.
	//  - I need the Authentication wiring.

	return nil
}

func Index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Index")
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	// TODO: if is authenticated redirect -> to index "/"
	io.WriteString(w, "You have to Login")
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

func Redirect(w http.ResponseWriter, loc string) {
	w.Header().Set("Location", loc)
	w.WriteHeader(http.StatusSeeOther)
}

func setupReload(router chi.Router) {
	// TODO: Implement Hot Reload
}

func mid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
