package middleware

import (
	"context"
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/rest"
	"github.com/danielcosme/curious-ape/rest/middleware"
	"net/http"
)

func SetDay(a *application.App) middleware.HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			var day *entity.Day

			if key := r.Header.Get("X-APE-DATE"); key != "" {
				date, err := entity.ParseDate(key)
				if err != nil {
					rest.ErrBadRequest(rw, r, err.Error())
					return
				}

				day, err = a.DayGetByDate(date)
				if err != nil {
					rest.ErrNotFound(rw, r)
				}
			}

			r = r.Clone(context.WithValue(r.Context(), "day", day))
			next.ServeHTTP(rw, r)
		})
	}
}
