package main

import (
	"encoding/base64"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

func (a *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			a.unauthorizedResponse(rw, r)
			return
		}

		auth := strings.Split(header, " ")[1:]
		if len(auth) != 1 {
			a.badRequestResponse(rw, r, errors.New("client needs to provide credentials"))
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(auth[0])
		if err != nil {
			a.serverErrorResponse(rw, r, err)
			return
		}

		auth = strings.Split(string(decoded), ":")
		if len(auth) != 2 {
			a.badRequestResponse(rw, r, errors.New("credentials need to be in username=password format"))
			return
		}
		usr := auth[0]
		pass := auth[1]

		user, err := a.models.Users.GetByEmail(usr)
		if err != nil {
			a.invalidCredentialsResponse(rw, r)
			return
		}

		isMatch, err := user.Password.IsMatch(pass)
		if err != nil {
			a.serverErrorResponse(rw, r, err)
			return
		}

		if !isMatch {
			a.invalidCredentialsResponse(rw, r)
			return
		}

		next.ServeHTTP(rw, r)
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
