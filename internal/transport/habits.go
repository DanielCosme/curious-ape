package transport

import (
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/view"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (t *Transport) newHabitLogPost(c echo.Context) error {
	success, err := strconv.ParseBool(c.QueryParam("success"))
	if err != nil {
		return errClientError(err)
	}
	date, err := core.DateFromISO8601(c.QueryParam("date"))
	if err != nil {
		return errClientError(err)
	}

	habit, err := t.App.HabitUpsert(core.NewHabitParams{
		Success:   success,
		Date:      date,
		HabitType: core.HabitType(c.QueryParam("category")),
		Origin:    core.OriginLogManual,
		Automated: false,
	})
	if err != nil {
		return errServer(err)
	}
	day, err := t.App.DayGetByID(habit.DayID)
	if err != nil {
		return errServer(err)
	}

	d := formatDays([]core.Day{day})[0]
	return t.RenderTempl(http.StatusOK, c, view.Day_Summary_Row(d))
}
