package web

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (h *Handler) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	h.App.Log.Error(errors.New(trace))

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (h *Handler) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (h *Handler) notFound(w http.ResponseWriter) {
	h.clientError(w, http.StatusNotFound)
}
