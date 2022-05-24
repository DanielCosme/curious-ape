package echo

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) FitbitDebug(c echo.Context) error {
	r, err := h.App.SleepDebug()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, r)
}
