package echo

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) DaysGetAll(c echo.Context) error {
	days, err := h.App.DaysGetAll()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, days)
}
