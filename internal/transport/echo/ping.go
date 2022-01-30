package echo

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) Ping(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
