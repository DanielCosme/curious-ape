package tucan

import (
	"net/http"
)

// TODO:
// - Add Middleware.
// 	- Rate limiter.
// 	- Panic recovery.
// 	- Request Logging.

//TODO:
// - Respond to requests.

// TODO: custom response controller.
// https://www.alexedwards.net/blog/how-to-use-the-http-responsecontroller-type

// Tucan is the main framework Instance.
type Tucan struct {
	middleware []MiddlewareFunc
}

// Context represents state of the current HTTP request.
type Context struct {
	Request  *http.Request
	Response *Response
}

// MiddlewareFunc defines a function to process middleware.
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

// HandlerFunc defines a function to serve HTTP requests.
type HandlerFunc func(c Context) error

func New() (t *Tucan) {
	return
}

func (t *Tucan) Use(ms ...MiddlewareFunc) {
	t.middleware = append(t.middleware, ms...)
}

func (t *Tucan) ServeHTTP(http.ResponseWriter, *http.Request) {}

type Response struct {
	http.ResponseWriter
}

func (r *Response) Header() http.Header {
	return nil
}

func (r *Response) Write([]byte) (int, error) {
	return 0, nil
}

func (r *Response) WriteHeader(statusCode int) {

}
