package transport

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (h *Transport) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	h.App.Log.Error(errors.New(trace))

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (h *Transport) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (h *Transport) notFound(w http.ResponseWriter) {
	h.clientError(w, http.StatusNotFound)
}
