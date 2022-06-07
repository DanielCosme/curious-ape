package middleware

import (
	"fmt"
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/rest"
	"github.com/danielcosme/curious-ape/sdk/errors"
	"net/http"
)

func RecoverPanic(a *application.App) rest.HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					rw.Header().Set("Connection", "close")
					rw.WriteHeader(http.StatusInternalServerError)
					rwPlus := rw.(*rest.ResponseWriterPlus)
					rwPlus.Err = errors.NewFatal(fmt.Sprintf("%v", err))
				}
			}()

			next.ServeHTTP(rw, r)
		})
	}
}
