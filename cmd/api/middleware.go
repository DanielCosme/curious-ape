package main

import (
	"expvar"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

func (a *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if a.config.env == "development" {
			next.ServeHTTP(rw, r)
			return
		}

		usr, pass, ok := r.BasicAuth()
		if ok {


			user, err := a.models.Users.GetByEmail(usr)
			if err != nil {
				a.serverErrorResponse(rw, r, err)
				return
			}

			isMatch, err := user.Password.IsMatch(pass)
			if err != nil {
				a.serverErrorResponse(rw, r, err)
				return
			}

			if isMatch {
				next.ServeHTTP(rw, r)
			}
		}

		a.invalidCredentialsResponse(rw, r)
	})
}

func (a *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic
		// as Go unwinds the stack).
		defer func() {
			if err := recover(); err != nil {
				rw.Header().Set("Connection", "close")
				a.serverErrorResponse(rw, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(rw, r)
	})
}

func (a *application) rateLimit(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(a.config.limiter.rps), a.config.limiter.burst)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.config.limiter.enabled {
			if !limiter.Allow() {
				a.rateLimitExceededResponse(w, r)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (a *application) metrics(next http.Handler) http.Handler {
	totalRequestsReceived := expvar.NewInt("total_requests_received")
	totalResponsesSent := expvar.NewInt("total_responses_sent")
	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_μs")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		totalRequestsReceived.Add(1)

		next.ServeHTTP(w, r)

		totalResponsesSent.Add(1)
		duration := time.Since(start).Microseconds()
		totalProcessingTimeMicroseconds.Add(duration)
	})
}
