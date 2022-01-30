package transport

import (
	"github.com/danielcosme/curious-ape/internal/transport/echo"
	"github.com/danielcosme/curious-ape/internal/transport/stdmux"
	"net/http"
)

type Router struct {
	StdMux http.Handler
	Echo   http.Handler
}

func (a *Transport) Routes() http.Handler {
	return &Router{
		StdMux: stdmux.Routes(a.App),
		Echo:   echo.Routes(a.App),
	}
}

func (router *Router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// we wrap the ServeMux in case in the future we want different handlers (echo, chi, etc) for different versions
	// 		of the API, v1, v2, etc...

	// router.StdMux.ServeHTTP(rw, r)
	router.Echo.ServeHTTP(rw, r)
}
