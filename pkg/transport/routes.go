package transport

import (
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
	// e.Use(t.midSecureHeaders)

	e.StaticFS("/static", echo.MustSubFS(web.Files, "static"))

	// login := e.Group("/login", t.midLoadAndSaveCookie)
	// {
	// 	login.GET("", t.loginForm)
	// 	login.POST("", t.loginPost)
	// }

	// p := e.Group("/" /*t.midLoadAndSaveCookie*/)
	{
		// p.GET("", t.home)
		// p.Use(t.midAuthenticateFromSession)
		// p.Use(t.midRequireAuth)

		// p.POST("logout", t.logout)
	}

	api := e.Group("/api/v1")
	{
		api.GET("", t.home)
	}

	// TODO make this endpoint protected.
	// e.GET("api/oauth2/:provider/success", t.oauth2Success)

	// In case I need a custom error Handler.
	// e.HTTPErrorHandler
	return e
}
