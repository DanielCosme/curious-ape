package httprest

import (
	"github.com/danielcosme/curious-ape/rest/middleware"
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
	// 		of the API, v1, v2, etc...
	mux.V1.ServeHTTP(rw, r)
}

func Routes(a *API) http.Handler {
	mux := http.NewServeMux()
	md := middleware.New()

	md.Use(middleware.LogRequest)

	mux.HandleFunc("/ping", a.Ping)
	mux.HandleFunc("/habits/", a.Habits)
	mux.HandleFunc("/", a.NotFound)

	return md.Commit(mux)
}
