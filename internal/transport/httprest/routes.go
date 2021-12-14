package httprest

import (
	"net/http"
)

type APIRouter struct {
	V1 http.Handler
}

func (a *API) Routes() http.Handler {
	return &APIRouter{
		V1: Routes(a),
	}
}

func (mux *APIRouter) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// we wrap the ServeMux in case in the future we want different handlers (echo, chi, etc) for different versions
	// 		of the API
	mux.V1.ServeHTTP(rw, r)
}

func Routes(a *API) http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("/", a.Ping)

	return r
}
