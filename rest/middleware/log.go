package middleware

import (
	"log"
	"net/http"
	"time"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		next.ServeHTTP(w, r)

		milis := time.Since(t).Milliseconds()
		log.Printf("%s - %s - %d Miliseconds", r.Method, r.URL.Path, milis)
	})
}
