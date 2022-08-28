package routes

import (
	"github.com/danielcosme/rest"
	"net/http"
)

func (h *Handler) Ping(rw http.ResponseWriter, r *http.Request) {
	rest.JSONStatusOk(rw, nil)
}
