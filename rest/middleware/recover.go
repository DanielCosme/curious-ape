package middleware

import (
	"github.com/danielcosme/curious-ape/rest"
	"log"
	"net/http"
)

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				rw.Header().Set(rest.HeaderConnection, "close")
				log.Printf("Panic: %s", err)
				rest.ErrInternalServer(rw, r)
			}
		}()

		next.ServeHTTP(rw, r)
	})
}

func Misc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			rest.ErrInternalServer(rw, r)
			return
		}

		next.ServeHTTP(rw, r)
	})
}
