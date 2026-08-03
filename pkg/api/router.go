package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/stephenafamo/bob"

	"danicos.dev/daniel/curious-ape/pkg/config"
	"danicos.dev/daniel/curious-ape/pkg/services/day"
	"danicos.dev/daniel/curious-ape/pkg/services/deadline"
	"danicos.dev/daniel/curious-ape/pkg/services/fitnesslog"
	"danicos.dev/daniel/curious-ape/pkg/services/habit"
	"danicos.dev/daniel/curious-ape/pkg/services/integration"
	"danicos.dev/daniel/curious-ape/pkg/services/sleeplog"
	"danicos.dev/daniel/curious-ape/pkg/services/user"
	"danicos.dev/daniel/curious-ape/pkg/services/worklog"
	"danicos.dev/daniel/curious-ape/pkg/ui"
	"danicos.dev/daniel/curious-ape/web/resources"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	sessionManager *scs.SessionManager
	Day            *day.Handler
	Habit          *habit.Handler
	Integration    *integration.Handler
	Worklog        *worklog.Handler
	User           *user.Handler
	Fitness        *fitnesslog.Handler
	Sleep          *sleeplog.Handler
	Deadline       *deadline.Handler
}

func NewHandlers(
	SessionManager *scs.SessionManager,
	Day *day.Service,
	Habit *habit.Service,
	Integration *integration.Service,
	Worklog *worklog.Service,
	User *user.Service,
	Fitness *fitnesslog.Service,
	Sleep *sleeplog.Service,
	Deadline *deadline.Service,
) Handlers {
	return Handlers{
		sessionManager: SessionManager,
		Day:            day.NewHandler(Day),
		Habit:          habit.NewHandler(Habit),
		Integration:    integration.NewHandler(Integration),
		Worklog:        worklog.NewHandler(Worklog),
		User:           user.NewHandler(User, SessionManager),
		Fitness:        fitnesslog.NewHandler(Fitness),
		Sleep:          sleeplog.NewHandler(Sleep),
		Deadline:       deadline.NewHandler(Deadline),
	}
}

func SetupRouter(ctx context.Context, handlers Handlers, r chi.Router, db bob.DB) (err error) {
	if config.IsDev() {
		setupReload(r)
	}

	r.Handle("/static/*", resources.Handler())

	r.Group(func(r chi.Router) {
		r.Use(handlers.sessionManager.LoadAndSave)
		r.Use(user.AuthenticateFromSession(handlers.sessionManager, db))

		r.Group(func(r chi.Router) {
			integration.SetupOauthRoutes(r, handlers.Integration)
		})

		r.Use(SetState)

		err = user.SetupRoutes(r, handlers.User)

		r.Group(func(r chi.Router) {
			r.Use(user.RequireAuthentication)

			day.SetupRoutes(r, handlers.Day)
			habit.SetupRoutes(r, handlers.Habit)
			integration.SetupRoutes(r, handlers.Integration)
			worklog.SetupRoutes(r, handlers.Worklog)
			fitnesslog.SetupRoutes(r, handlers.Fitness)
			sleeplog.SetupRoutes(r, handlers.Sleep)
			deadline.SetupRoutes(r, handlers.Deadline)
		})
	})

	if err != nil {
		return fmt.Errorf("error setting up routes: %w", err)
	}
	for _, route := range r.Routes() {
		slog.Debug(fmt.Sprintf("Service registered: %s", route.Pattern))
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
