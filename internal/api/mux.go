package api

import (
	"github.com/danielcosme/curious-ape/internal/api/router"
	"net/http"
)

type Mux struct {
	StdMux http.Handler
	// Echo   http.Handler
}

func (t *Transport) Routes() http.Handler {
	return &Mux{
		StdMux: router.Routes(t.App),
		// Echo:   echo.Routes(a.App),
	}
}

func (router *Mux) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// we wrap the ServeMux in case in the future we want different handlers (echo, chi, etc) for different versions
	// 		of the API, v1, v2, etc...

	router.StdMux.ServeHTTP(rw, r)
	// router.Echo.ServeHTTP(rw, r)
}
