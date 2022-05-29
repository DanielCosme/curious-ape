package rest

import (
	"context"
	"net/http"
)

func ErrResponse(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	r = r.Clone(context.WithValue(r.Context(), "error", payload))
	switch payload.(type) {
	case error:
		payload = payload.(error).Error()
	}

	JSON(w, code, &envelope{"error": payload})
}

func ErrBadRequest(rw http.ResponseWriter, r *http.Request, payload interface{}) {
	ErrResponse(rw, r, http.StatusBadRequest, payload)
}

func ErrInternalServer(rw http.ResponseWriter, r *http.Request) {
	ErrResponse(rw, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func ErrNotFound(rw http.ResponseWriter, r *http.Request) {
	ErrResponse(rw, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

func ErrNotAllowed(rw http.ResponseWriter, r *http.Request) {
	ErrResponse(rw, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
}
