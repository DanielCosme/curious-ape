package router

import (
	"github.com/danielcosme/curious-ape/rest"
	"net/http"
)

func (a Handler) Ping(rw http.ResponseWriter, r *http.Request) {
	rest.JSONStatusOk(rw, nil)
}

func (a Handler) NotFound(rw http.ResponseWriter, r *http.Request) {
	rest.ErrNotFound(rw, r)
}
