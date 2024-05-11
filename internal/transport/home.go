package transport

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"

	"github.com/danielcosme/curious-ape/internal/core"
)

func (t *Transport) home(c echo.Context) error {
	ds, err := t.App.DaysCurMonth()
	if err != nil {
		return errServer(err)
	}
	data := t.newTemplateData(c.Request())
	data.Days = formatDays(ds)
	return c.Render(http.StatusOK, "home.gohtml", data)
}

type dayContainer struct {
	Date    time.Time
	Wake    *core.Habit
	Fitness *core.Habit
	Work    *core.Habit
	Eat     *core.Habit
}

func formatDays(ds []*core.Day) []dayContainer {
	var res []dayContainer
	for _, d := range ds {
		dc := dayContainer{Date: d.Date.Time()}
		for _, h := range d.Habits {
			switch h.Category.Type {
			case core.HabitTypeWakeUp:
				dc.Wake = &h
			case core.HabitTypeExercise:
				dc.Fitness = &h
			case core.HabitTypeDeepWork:
				dc.Work = &h
			case core.HabitTypeEatHealthy:
				dc.Eat = &h
			}
		}
		dc.Wake = replace(dc.Wake)
		dc.Fitness = replace(dc.Fitness)
		dc.Work = replace(dc.Work)
		dc.Eat = replace(dc.Eat)
		res = append(res, dc)
	}
	return res
}

func replace(h *core.Habit) *core.Habit {
	if h == nil {
		return core.NewHabit()
	}
	return h
}
