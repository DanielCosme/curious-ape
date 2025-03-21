package transport

type ContextKey string

const (
	ctxKeyIsAuthenticated     ContextKey = "isAuthenticated"
	ctxKeyAuthenticatedUserID ContextKey = "authenticatedUserID"
	ctxUser                   ContextKey = "user"
)
