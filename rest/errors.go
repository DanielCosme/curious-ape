package rest

import (
	"fmt"
	"net/http"
)

func ErrResponse(rw http.ResponseWriter, code int, payload interface{}) {
	var p interface{}

	rwPlus := rw.(*ResponseWriterPlus)
	switch payload.(type) {
	case error:
		rwPlus.Err = payload.(error)
		p = rwPlus.Err.Error()
	default:
		p = payload
		rwPlus.Err = fmt.Errorf("%v", payload)
	}

	JSON(rw, code, &envelope{"error": p})
}

func ErrBadRequest(rw http.ResponseWriter, payload interface{}) {
	ErrResponse(rw, http.StatusBadRequest, payload)
}

func ErrInternalServer(rw http.ResponseWriter) {
	ErrResponse(rw, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func ErrNotFound(rw http.ResponseWriter) {
	ErrResponse(rw, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

func ErrNotAllowed(rw http.ResponseWriter) {
	ErrResponse(rw, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
}
