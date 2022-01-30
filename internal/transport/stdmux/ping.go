package stdmux

import (
	"github.com/danielcosme/curious-ape/rest"
	"net/http"
)

func (a Handler) Ping(rw http.ResponseWriter, r *http.Request) {
	rest.JSONStatusOk(rw, nil)
}
