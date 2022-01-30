package stdmux

import (
	"github.com/danielcosme/curious-ape/rest"
	"net/http"
)

func (a Handler) NotFound(rw http.ResponseWriter, r *http.Request) {
	rest.ErrNotFound(rw)
}
