package rest

import (
	"github.com/danielcosme/curious-ape/sdk/errors"
	"net/http"
)

type ResponseWriterPlus struct {
	http.ResponseWriter
	status int
	Err    error
}

func NewResponseWriterPlus(rw http.ResponseWriter) *ResponseWriterPlus {
	return &ResponseWriterPlus{ResponseWriter: rw}
}

func (rw *ResponseWriterPlus) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *ResponseWriterPlus) Status() int {
	if rw.status == 0 {
		return http.StatusOK
	}
	if rw.Err == nil && rw.status == http.StatusMethodNotAllowed {
		// TODO handle this properly by setting a method not allowed handler on chi routes
		rw.Err = errors.New("unknown error")
	}
	return rw.status
}
