package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stephenafamo/bob"

	"danicos.dev/daniel/curious-ape/pkg/config"
	"danicos.dev/daniel/curious-ape/pkg/services/day"
	"danicos.dev/daniel/curious-ape/pkg/services/habit"
	"danicos.dev/daniel/curious-ape/pkg/services/integration"
	"danicos.dev/daniel/curious-ape/pkg/services/user"
	"danicos.dev/daniel/curious-ape/pkg/services/worklog"
	"danicos.dev/daniel/curious-ape/pkg/ui"
	"danicos.dev/daniel/curious-ape/web/resources"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

type Handllers struct {
	Day         *day.Handler
	Habit       *habit.Handler
	Integration *integration.Handler
	Worklog     *worklog.Handler
	User        *user.Handler
}

func SetupRouter(ctx context.Context, handlers Handllers, cfg *config.Config, r chi.Router, sessionManager *scs.SessionManager, db bob.DB) (err error) {
	if cfg.Environment == config.Dev {
		setupReload(r)
	}

	r.Handle("/static/*", resources.Handler())

	r.Group(func(r chi.Router) {
		r.Use(sessionManager.LoadAndSave)
		r.Use(user.AuthenticateFromSession(sessionManager, db))
		r.Use(SetState)

		err = user.SetupRoutes(r, handlers.User)

		r.Group(func(r chi.Router) {
			r.Use(user.RequireAuthentication)

			day.SetupRoutes(r, handlers.Day)
			habit.SetupRoutes(r, handlers.Habit)
			integration.SetupRoutes(r, handlers.Integration)
			worklog.SetupRoutes(r, handlers.Worklog)
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
			CurrentPath:     r.URL.Path,
		}
		ctx := ui.StateWithContextUI(r.Context(), &state)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func setupReload(router chi.Router) {
	// TODO: Implement Hot Reload for development
}
