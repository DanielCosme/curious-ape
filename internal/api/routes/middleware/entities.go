package middleware

import (
	"context"
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/rest"
	"github.com/danielcosme/curious-ape/sdk/errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func SetHabit(a *application.App) rest.HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			id, err := parseID(chi.URLParam(r, "habitID"))
			if err != nil {
				rest.ErrInternalServer(rw)
				return
			}

			habit, err := a.HabitGetByID(id)
			if err != nil {
				rest.ErrNotFound(rw)
				return
			}

			r = r.Clone(context.WithValue(r.Context(), "habit", habit))
			next.ServeHTTP(rw, r)
		})
	}
}

func SetSleepLog(a *application.App) rest.HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			id, err := parseID(chi.URLParam(r, "id"))
			if err != nil {
				rest.ErrInternalServer(rw)
				return
			}

			sl, err := a.GetSleepLog(entity.SleepLogFilter{ID: []int{id}})
			if err != nil {
				rest.ErrNotFound(rw)
				return
			}

			r = r.Clone(context.WithValue(r.Context(), "sleepLog", sl))
			next.ServeHTTP(rw, r)
		})
	}
}

func SetFitnessLog(a *application.App) rest.HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			id, err := parseID(chi.URLParam(r, "id"))
			if err != nil {
				rest.ErrInternalServer(rw)
				return
			}

			fl, err := a.GetFitnessLog(entity.FitnessLogFilter{ID: []int{id}})
			if err != nil {
				rest.ErrNotFound(rw)
				return
			}

			r = r.Clone(context.WithValue(r.Context(), "fitnessLog", fl))
			next.ServeHTTP(rw, r)
		})
	}
}

func parseID(idStr string) (int, error) {
	if idStr != "" {
		return strconv.Atoi(idStr)
	}
	return 0, errors.New("id param not found")

}
