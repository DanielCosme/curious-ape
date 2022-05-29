package router

import (
	"github.com/danielcosme/curious-ape/rest"
	"github.com/danielcosme/curious-ape/sdk/errors"
	"net/http"
)

func (a Handler) Ping(rw http.ResponseWriter, r *http.Request) {
	rest.ErrResponse(rw, http.StatusInternalServerError, errors.NewFatal("something went wrong"))
	return
	rest.JSONStatusOk(rw, nil)
}

