package main

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

func (a *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic
		// as Go unwinds the stack).
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				a.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (a *application) rateLimit(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			a.rateLimitExceededResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
