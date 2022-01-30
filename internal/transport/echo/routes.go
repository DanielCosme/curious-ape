package echo

import (
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func Routes(a *application.App) http.Handler {
	h := Handler{App: a}
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/ping", h.Ping)
	e.GET("/habits", h.HabitsGetAll)
	e.GET("/habits/:id", h.HabitsGetAll)

	return e
}
