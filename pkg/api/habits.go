package api

import (
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (api *API) updateHabit(c echo.Context) error {
	date, err := core.DateFromISO8601(c.QueryParam("date"))
	if err != nil {
		return errClientError(err)
	}
	habitType := c.QueryParam("type")
	habitState := c.QueryParam("state")

	_, err = api.App.HabitUpsert(date, core.HabitTypeFromString(habitType), core.HabitState(habitState))
	if err != nil {
		return errServer(err)
	}
	d, err := api.App.DayGetOrCreate(date)
	if err != nil {
		return errServer(err)
	}
	return c.JSON(http.StatusOK, dayDBToSummary(d))
}
