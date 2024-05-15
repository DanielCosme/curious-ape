package transport

import (
	"github.com/danielcosme/curious-ape/web"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func EchoRoutes(t *Transport) http.Handler {
	e := echo.New()
	e.Renderer = t

	e.Use(middleware.RequestLoggerWithConfig(midSlogConfig(t)))
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(t.midSecureHeaders)

	e.StaticFS("/static", echo.MustSubFS(web.Files, "static"))

	login := e.Group("/login", t.midLoadAndSaveCookie)
	login.GET("", t.loginForm)
	login.POST("", t.loginPost)

	p := e.Group("/", t.midLoadAndSaveCookie)
	{
		p.Use(t.midAuthenticateFromSession)
		p.Use(t.midRequireAuth)

		p.GET("", t.home)
		p.POST("logout", t.logout)

		p.POST("habit/log", t.newHabitLogPost)

		p.GET("integrations", t.integrationsForm)
	}

	e.GET("api/oauth2/:provider/success", t.oauth2Success)

	// In case I need a custom error Handler.
	// e.HTTPErrorHandler
	return e
}
