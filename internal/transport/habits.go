package transport

import (
	"github.com/danielcosme/curious-ape/internal/core"
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
	categoryID, err := strconv.Atoi(c.QueryParam("category"))
	if err != nil {
		return errClientError(err)
	}

	habit, err := t.App.HabitUpsert(core.HabitParams{
		Success:    success,
		Date:       date,
		CategoryID: int32(categoryID),
		Origin:     core.WebUI,
		Automated:  false,
	})
	if err != nil {
		return errServer(err)
	}
	day, err := t.App.DayGetByID(habit.DayID)
	if err != nil {
		return errServer(err)
	}

	td := t.newTemplateData(c.Request())
	td.Day = &formatDays([]core.Day{day})[0]
	return c.Render(http.StatusOK, partial("day_row.gohtml"), td.Day)
}
