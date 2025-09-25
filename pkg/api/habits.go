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

	_, err = api.App.HabitUpsert(core.UpsertHabitParams{
		Date:  date,
		Type:  core.HabitTypeFromString(c.QueryParam("type")),
		State: core.HabitState(c.QueryParam("state"))})
	if err != nil {
		return errServer(err)
	}
	d, err := api.App.DayGetOrCreate(date)
	if err != nil {
		return errServer(err)
	}
	return c.JSON(http.StatusOK, d)
}
