package echo

import (
	"github.com/danielcosme/curious-ape/sdk/errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) Ping(c echo.Context) error {
	return  errors.NewFatal("This is a fatal error")
	return c.NoContent(http.StatusOK)
}
