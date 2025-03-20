package transport

import (
	"github.com/danielcosme/curious-ape/pkg/application"
	"github.com/danielcosme/curious-ape/web"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func EchoRoutes(t *Transport) http.Handler {
	e := echo.New()

	e.Use(middleware.RequestLoggerWithConfig(midSlogConfig(t)))
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(t.midSecureHeaders)

	e.StaticFS("/static", echo.MustSubFS(web.Files, "static"))
	// e.StaticFS("/", echo.MustSubFS(web.Files, "dist"))

	login := e.Group("/api/v1/login", t.midLoadAndSaveCookie)
	{
		login.POST("", t.loginPost)
	}

	api := e.Group("/api/v1", t.midLoadAndSaveCookie)
	{
		api.Use(t.midAuthenticateFromSession)
		api.Use(t.midRequireAuth)

		api.GET("", t.home)

		api.GET("/integrations", t.integrationsGetAll)
		api.GET("/integrations/:provider", t.integrationsGet)
		api.POST("/sync", t.syncDay)

		// p.POST("logout", t.logout)
	}
	e.GET("api/oauth2/:provider/success", t.oauth2Success)

	if t.App.Env == application.Dev {
		debug := e.Group("/api/v1/debug")
		{
			debug.GET("", t.home)
		}
	}

	// In case I need a custom error Handler.
	// e.HTTPErrorHandler
	return e
}
