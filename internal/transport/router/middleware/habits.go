package middleware

import (
	"context"
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/rest"
	"net/http"
	"strconv"
)

func SetHabit(a *application.App) rest.HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			var habit *entity.Habit

			if idStr := r.Form.Get("habit_id"); idStr != "" {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					rest.ErrInternalServer(rw, r)
					return
				}

				habit, err = a.HabitGetByID(id)
				if err != nil {
					rest.ErrNotFound(rw, r)
					return
				}
			}

			r = r.Clone(context.WithValue(r.Context(), "habit", habit))
			next.ServeHTTP(rw, r)
		})
	}
}
