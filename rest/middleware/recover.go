package middleware

import "net/http"

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				rw.Header().Set("Connection", "close")
				res := "Internal server error"
				rw.Write([]byte(res))
				rw.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(rw, r)
	})
}
