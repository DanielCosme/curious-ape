package dove

import "net/http"

type endpoint struct {
	Path       string
	Handlers   map[string]HandlerFunc
	middleware []MiddlewareFunc
}

func Endpoint(path string) *endpoint {
	e := endpoint{
		Handlers: map[string]HandlerFunc{},
	}
	return &e
}

func (e *endpoint) GET(fn HandlerFunc) *endpoint {
	if fn != nil {
		fn = e.addMiddleware(fn)
		e.Handlers[http.MethodGet] = fn
	}
	return e
}

func (e *endpoint) POST(fn HandlerFunc) *endpoint {
	if fn != nil {
		fn = e.addMiddleware(fn)
		e.Handlers[http.MethodPost] = fn
	}
	return e
}

func (e *endpoint) PUT(fn HandlerFunc) *endpoint {
	if fn != nil {
		fn = e.addMiddleware(fn)
		e.Handlers[http.MethodPut] = fn
	}
	return e
}

func (e *endpoint) PATCH(fn HandlerFunc) *endpoint {
	if fn != nil {
		fn = e.addMiddleware(fn)
		e.Handlers[http.MethodPatch] = fn
	}
	return e
}

func (e *endpoint) DELETE(fn HandlerFunc) *endpoint {
	if fn != nil {
		fn = e.addMiddleware(fn)
		e.Handlers[http.MethodDelete] = fn
	}
	return e
}

func (e *endpoint) addMiddleware(fn HandlerFunc) HandlerFunc {
	for i := len(e.middleware) - 1; i >= 0; i-- {
		fn = e.middleware[i](fn)
	}
	return fn
}
