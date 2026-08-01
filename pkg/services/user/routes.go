package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ContextKey string

const (
	ctxKeyIsAuthenticated     ContextKey = "isAuthenticated"
	ctxKeyAuthenticatedUserID ContextKey = "authenticatedUserID"
	ctxUser                   ContextKey = "user"
)

func SetupRoutes(r chi.Router, handler *Handler) error {
	r.Route("/login", func(r chi.Router) {
		r.Get("/", handler.LoginPage)
		r.Post("/", handler.LoginPost)
		r.Delete("/", handler.Logout)
	})
	return nil
}

func IsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(ctxKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
