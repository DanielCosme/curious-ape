package transport

import (
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/danielcosme/curious-ape/internal/application"
	entity2 "github.com/danielcosme/curious-ape/internal/entity"

	"github.com/danielcosme/curious-ape/internal/validator"
)

type newHabitForm struct {
	Date         time.Time
	CategoryCode string
	Success      bool
	Origin       entity2.DataSource
	Note         string
	IsAutomated  bool
	validator.Validator
}

func (t *Transport) habit(c echo.Context) error {
	habit := c.Get("habit").(*entity2.Habit)
	data := t.newTemplateData(c.Request())
	data.Habit = habit
	return c.Render(http.StatusOK, "view.gohtml", data)
}

func (t *Transport) newHabitForm(c echo.Context) error {
	data := t.newTemplateData(c.Request())
	data.Form = newHabitForm{}
	return c.Render(http.StatusOK, "new_habit.gohtml", data)
}

func (t *Transport) newHabitLogPost(c echo.Context) error {
	slog.Info("HERE")
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
