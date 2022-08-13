package rest

import (
	"net/http"
)

func MiddlewareRecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				rw.Header().Set(HeaderConnection, "close")
				ErrInternalServer(rw)
			}
		}()

		next.ServeHTTP(rw, r)
	})
}

func MiddlewareParseForm(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			ErrInternalServer(rw)
			return
		}

		next.ServeHTTP(rw, r)
	})
}

func AllowCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// TODO implement trusted origins
		// rw.Header().Add("Vary", "Origin")
		// origin := r.Header.Get("Origin")

		rw.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(rw, r)
	})
}
