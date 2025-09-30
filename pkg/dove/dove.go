package dove

import (
	"net/http"

	"github.com/danielcosme/curious-ape/pkg/oak"
)

// TODO:
// - Add Middleware.
// 	- Rate limiter.

// TODO: custom response controller.
// https://www.alexedwards.net/blog/how-to-use-the-http-responsecontroller-type

// Dove is the main framework Instance.
type Dove struct {
	middleware []MiddlewareFunc
	routes     map[string]*endpoint
}

func New() *Dove {
	d := &Dove{
		routes: map[string]*endpoint{},
	}
	return d
}

func (d *Dove) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	endpoint, ok := d.routes[r.URL.Path]
	if !ok {
		oak.Info("not found: " + r.URL.Path)
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	handler, ok := endpoint.Handlers[r.Method]
	if !ok {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	for i := len(d.middleware) - 1; i >= 0; i-- {
		handler = d.middleware[i](handler)
	}

	doveCtx := Context{
		Req: r,
		Res: rw,
	}
	// TODO: Implement out error handling.
	_ = handler(doveCtx)
}

func (d *Dove) Endpoint(path string) *endpoint {
	e := Endpoint(path)
	d.routes[path] = e
	return e
}

func (d *Dove) Use(ms ...MiddlewareFunc) {
	d.middleware = append(d.middleware, ms...)
}
