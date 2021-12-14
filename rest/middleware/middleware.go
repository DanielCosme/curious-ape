package middleware

import "net/http"

type Middleware struct {
	middlewares []HTTPMiddleware
}

type HTTPMiddleware func(http.Handler) http.Handler

func New(mds ...HTTPMiddleware) *Middleware {
	return &Middleware{
		middlewares: mds,
	}
}

func (m *Middleware) Use(middleware HTTPMiddleware) {
	m.middlewares = append(m.middlewares, middleware)
}

func (m *Middleware) Commit(h http.Handler) http.Handler {
	for i := len(m.middlewares) -1; i >= 0 ; i-- {
		h = m.middlewares[i](h)
	}

	return h
}

func(m *Middleware) Then(f http.HandlerFunc) http.Handler{
	return m.Commit(f)
}
