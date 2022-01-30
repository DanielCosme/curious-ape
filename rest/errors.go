package rest

import "net/http"

func ErrResponse(w http.ResponseWriter, code int, payload interface{}) {
	JSON(w, code, Payload("error", payload))
}

func ErrBadRequest(w http.ResponseWriter, payload interface{}) {
	ErrResponse(w, http.StatusBadRequest, payload)
}

func ErrInternalServer(w http.ResponseWriter) {
	s := http.StatusInternalServerError
	ErrResponse(w, s, http.StatusText(s))
}

func ErrNotFound(w http.ResponseWriter) {
	s := http.StatusNotFound
	ErrResponse(w, s, http.StatusText(s))
}

func ErrMethodNotSupported(w http.ResponseWriter) {
	s := http.StatusMethodNotAllowed
	ErrResponse(w, s, http.StatusText(s))
}
