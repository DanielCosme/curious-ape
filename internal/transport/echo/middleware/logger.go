package middleware

import (
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/labstack/echo/v4"
)

func Logger(a *application.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			return next(c)
		}
	}
}

func Recover(a *application.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			return next(c)
		}
	}
}
