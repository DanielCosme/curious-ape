package router

import (
	"github.com/danielcosme/curious-ape/rest"
	"net/http"
)

func (h *Handler) Ping(rw http.ResponseWriter, r *http.Request) {
	rest.JSONStatusOk(rw, nil)
}
