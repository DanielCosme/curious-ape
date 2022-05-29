package router

import (
	"github.com/danielcosme/curious-ape/rest"
	"net/http"
)

func (a Handler) NotFound(rw http.ResponseWriter, r *http.Request) {
	rest.ErrNotFound(rw, r)
}

func JsonCheckError(rw http.ResponseWriter, r *http.Request, status int, data interface{}, err error) {
	if err != nil {
		// TODO separate not-found errors from internal server errors
		rest.ErrResponse(rw, r, http.StatusInternalServerError, err)
	} else {
		rest.JSON(rw, status, data)
	}
}
