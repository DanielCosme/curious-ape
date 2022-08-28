package middleware

import (
	"context"
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/rest"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func SetDay(a *application.App) rest.HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			var day *entity.Day

			if key := chi.URLParam(r, "date"); key != "" {
				date, err := entity.ParseDate(key)
				if err != nil {
					rest.ErrBadRequest(rw, err.Error())
					return
				}

				day, err = a.DayGetByDate(date)
				if err != nil {
					rest.ErrResponse(rw, http.StatusNotFound, err)
					return
				}
			}

			r = r.Clone(context.WithValue(r.Context(), "day", day))
			next.ServeHTTP(rw, r)
		})
	}
}
