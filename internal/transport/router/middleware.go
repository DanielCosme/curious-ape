package router

//
// import (
// 	"expvar"
// 	"fmt"
// 	"net/http"
// 	"time"
//
// 	"golang.org/x/time/rate"
// )
//
// func (a *main.application) authenticate(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
// 		if a.config.env == "development" {
// 			next.ServeHTTP(rw, r)
// 			return
// 		}
//
// 		usr, pass, ok := r.BasicAuth()
// 		if ok {
//
// 			user, err := a.models.Users.GetByEmail(usr)
// 			if err != nil {
// 				a.serverErrorResponse(rw, r, err)
// 				return
// 			}
//
// 			if string(user.Password.Hash) == pass {
// 				fmt.Print("Authentication ok")
// 				next.ServeHTTP(rw, r)
// 				return
// 	 		}
// 		}
//
// 		a.invalidCredentialsResponse(rw, r)
// 	})
// }
//
//
// func (a *main.application) rateLimit(next http.Handler) http.Handler {
// 	limiter := rate.NewLimiter(rate.Limit(a.config.limiter.rps), a.config.limiter.burst)
//
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if a.config.limiter.enabled {
// 			if !limiter.Allow() {
// 				a.rateLimitExceededResponse(w, r)
// 				return
// 			}
// 		}
//
// 		next.ServeHTTP(w, r)
// 	})
// }
//
// func (a *main.application) metrics(next http.Handler) http.Handler {
// 	totalRequestsReceived := expvar.NewInt("total_requests_received")
// 	totalResponsesSent := expvar.NewInt("total_responses_sent")
// 	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_μs")
//
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		totalRequestsReceived.Add(1)
//
// 		next.ServeHTTP(w, r)
//
// 		totalResponsesSent.Add(1)
// 		duration := time.Since(start).Microseconds()
// 		totalProcessingTimeMicroseconds.Add(duration)
// 	})
// }
//
