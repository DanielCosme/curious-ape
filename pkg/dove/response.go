package dove

import (
	"net/http"

	"git.danicos.dev/daniel/curious-ape/pkg/oak"
)

type Response struct {
	Writer      http.ResponseWriter
	Logger      *oak.Oak
	StatusCode  int
	Size        uint64
	Commited    bool
	beforeFuncs []func()
	// afterFuncs  []func()
}

func NewResponse(w http.ResponseWriter, l *oak.Oak) *Response {
	return &Response{
		Writer: w,
		Logger: l,
	}
}

func (r *Response) Before(fn func()) {
	r.beforeFuncs = append(r.beforeFuncs, fn)
}

// Header returns the header map for the writer that will be sent by
// WriteHeader. Changing the header after a call to WriteHeader (or Write) has
// no effect unless the modified headers were declared as trailers by setting
// the "Trailer" header before the call to WriteHeader (see example)
// To suppress implicit response headers, set their value to nil.
// Example: https://golang.org/pkg/net/http/#example_ResponseWriter_trailers
func (r *Response) Header() http.Header {
	return r.Writer.Header()
}

// WriteHeader sends an HTTP response header with status code. If WriteHeader is
// not called explicitly, the first call to Write will trigger an implicit
// WriteHeader(http.StatusOK). Thus explicit calls to WriteHeader are mainly
// used to send error codes.
func (r *Response) WriteHeader(statusCode int) {
	if r.Commited {
		r.Logger.Error("response already commited")
		return
	}
	r.StatusCode = statusCode
	for _, fn := range r.beforeFuncs {
		fn()
	}
	r.Writer.WriteHeader(r.StatusCode)
	r.Commited = true
}

// Write writes the data to the connection as part of an HTTP reply.
func (r *Response) Write(body []byte) (int, error) {
	if !r.Commited {
		if r.StatusCode == 0 {
			r.StatusCode = http.StatusOK
		}
		r.WriteHeader(r.StatusCode)
	}

	n, err := r.Writer.Write(body)
	r.Size += uint64(n)
	return n, err
}
