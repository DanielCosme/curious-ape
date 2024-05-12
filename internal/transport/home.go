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
	Wake    core.Habit
	Fitness core.Habit
	Work    core.Habit
	Eat     core.Habit
}

func formatDays(ds []core.Day) []dayContainer {
	var res []dayContainer
	for _, d := range ds {
		dc := dayContainer{Date: d.Date.Time()}
		for _, h := range d.Habits {
			switch h.Category.Type {
			case core.HabitTypeWakeUp:
				dc.Wake = h
			case core.HabitTypeExercise:
				dc.Fitness = h
			case core.HabitTypeDeepWork:
				dc.Work = h
			case core.HabitTypeEatHealthy:
				dc.Eat = h
			}
		}
		dc.Wake = replace(dc.Wake)
		dc.Wake.Category.ID = 1
		dc.Fitness = replace(dc.Fitness)
		dc.Fitness.Category.ID = 2
		dc.Work = replace(dc.Work)
		dc.Work.Category.ID = 3
		dc.Eat = replace(dc.Eat)
		dc.Eat.Category.ID = 4
		res = append(res, dc)
	}
	return res
}

func replace(h core.Habit) core.Habit {
	if h.IsZero() {
		// TODO: replace this for something better.
		return core.NewHabit(core.NewDate(time.Now()), core.HabitCategory{}, nil)
	}
	return h
}
