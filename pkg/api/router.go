package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stephenafamo/bob"

	"danicos.dev/daniel/curious-ape/pkg/config"
	"danicos.dev/daniel/curious-ape/pkg/services/day"
	"danicos.dev/daniel/curious-ape/pkg/services/user"
	"danicos.dev/daniel/curious-ape/pkg/ui"
	"danicos.dev/daniel/curious-ape/web/resources"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(ctx context.Context, r chi.Router, sessionManager *scs.SessionManager, db bob.DB) (err error) {
	if config.Global.Environment == config.Dev {
		setupReload(r)
	}

	r.Handle("/static/*", resources.Handler())

	r.Group(func(r chi.Router) {
		r.Use(sessionManager.LoadAndSave)
		r.Use(user.AuthenticateFromSession(sessionManager, db))
		r.Use(SetState)

		err = user.SetupRoutes(r, db, sessionManager)

		r.Group(func(r chi.Router) {
			r.Use(user.RequireAuthentication)

			day.SetupRoutes(r, db)
		})
	})

	if err != nil {
		return fmt.Errorf("error setting up routes: %w", err)
	}
	return nil
}

func SetState(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state := ui.UIState{
			IsAuthenticated: user.IsAuthenticated(r),
			Version:         config.Version(),
		}
		ctx := ui.StateWithContextUI(r.Context(), &state)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func setupReload(router chi.Router) {
	// TODO: Implement Hot Reload for development
}
