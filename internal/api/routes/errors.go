package routes

import (
	"github.com/danielcosme/rest"
	"net/http"
)

func (h *Handler) NotFound(rw http.ResponseWriter, r *http.Request) {
	rest.ErrResponse(rw, http.StatusNotFound, "")
}

func (h *Handler) ErrInternalServer(rw http.ResponseWriter, msg string) {
	rest.ErrResponse(rw, http.StatusInternalServerError, msg)
}

func (h *Handler) ErrNotFound(rw http.ResponseWriter, msg string) {
	rest.ErrResponse(rw, http.StatusNotFound, msg)
}
