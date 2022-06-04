package router

import (
	"github.com/danielcosme/curious-ape/rest"
	"net/http"
)

func (a Handler) NotFound(rw http.ResponseWriter, r *http.Request) {
	rest.ErrResponse(rw, http.StatusNotFound, "")
}

func (a Handler) ErrInternalServer(rw http.ResponseWriter, msg string) {
	rest.ErrResponse(rw, http.StatusInternalServerError, msg)
}

func (a Handler) ErrNotFound(rw http.ResponseWriter, msg string) {
	rest.ErrResponse(rw, http.StatusNotFound, msg)
}
