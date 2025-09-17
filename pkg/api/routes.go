package api

import (
	"net/http"

	"github.com/danielcosme/curious-ape/views"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func EchoRoutes(a *API) http.Handler {
	e := echo.New()

	e.Use(middleware.RequestLoggerWithConfig(midSlogConfig(a)))
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(a.midSecureHeaders)

	e.StaticFS("/static", echo.MustSubFS(views.StaticFS, "static"))

	e.Use(a.midLoadAndSaveCookie)
	e.GET("/login", a.getLogin)
	e.POST("/login", a.loginPost)

	api := e.Group("/api/v1")
	{
		api.Use(a.midAuthenticateFromSession)
		api.Use(a.midRequireAuth)

		api.GET("", a.home)
		api.GET("/user", a.getUser)
		api.POST("/logout", a.logout)
		api.POST("/days/sync", a.syncDay)
		api.POST("/habits/update", a.updateHabit)
		api.GET("/integrations", a.integrationsGetAll)
		api.GET("/integrations/:provider", a.integrationsGet)
	}

	// In case I need a custom error Handler.
	// e.HTTPErrorHandler
	return e
}
