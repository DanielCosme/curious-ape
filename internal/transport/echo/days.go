package echo

import (
	"github.com/danielcosme/curious-ape/internal/transport/types"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) DaysGetAll(c echo.Context) error {
	days, err := h.App.DaysGetAll()
	if err != nil {
		return err
	}

	daysTransport := []*types.DayTransport{}
	for _, d := range days {
		daysTransport = append(daysTransport, types.DayToTransport(d))
	}

	return c.JSON(http.StatusOK, daysTransport)
}
