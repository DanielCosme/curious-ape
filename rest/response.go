package rest

import "net/http"

type ResponseWriterPlus struct {
	http.ResponseWriter
	status   int
	Err      error
	HasPanic bool
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
	return rw.status
}
