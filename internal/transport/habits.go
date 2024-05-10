package transport

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"

	"github.com/danielcosme/curious-ape/internal/application"
	entity2 "github.com/danielcosme/curious-ape/internal/entity"
)

func (t *Transport) newHabitLogPost(c echo.Context) error {
	success, err := strconv.ParseBool(c.QueryParam("success"))
	if err != nil {
		return errServer(err)
	}
	dt, err := time.Parse(time.DateOnly, c.QueryParam("date"))
	if err != nil {
		return errServer(err)
	}

	params := &application.NewHabitParams{
		Date:         dt,
		CategoryCode: c.QueryParam("category"),
		Success:      success,
		Origin:       entity2.Manual,
		IsAutomated:  false,
	}
	habit, err := t.App.HabitUpsert(params)
	if err != nil {
		return errServer(err)
	}
	day, err := t.App.DayGetByID(habit.DayID)
	if err != nil {
		return errServer(err)
	}

	td := t.newTemplateData(c.Request())
	td.Day = &formatDays([]*entity2.Day{day})[0]
	return c.Render(http.StatusOK, partial("day_row.gohtml"), td)
}
