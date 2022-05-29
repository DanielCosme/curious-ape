package middleware

import (
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/rest"
	"net/http"
)

func RecoverPanic(a *application.App) rest.HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {



			next.ServeHTTP(rw, r)
		})
	}
}
