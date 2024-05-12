package transport

type ContextKey string

const (
	ctxKeyIsAuthenticated     ContextKey = "isAuthenticated"
	ctxKeyAuthenticatedUserID ContextKey = "authenticatedUserID"
)

const (
	pageHome  = "home.gohtml"
	pageLogin = "login.gohtml"
)

const (
	partialDayRow = "p.day_row.gohtml"
)
