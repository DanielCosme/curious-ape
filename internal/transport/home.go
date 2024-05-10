package transport

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"

	entity2 "github.com/danielcosme/curious-ape/internal/entity"
)

func (t *Transport) home(c echo.Context) error {
	ds, err := t.App.DaysMonth()
	if err != nil {
		return errServer(err)
	}
	data := t.newTemplateData(c.Request())
	data.Days = formatDays(ds)
	return c.Render(http.StatusOK, "home.gohtml", data)
}

type dayContainer struct {
	Date    time.Time
	Wake    *entity2.Habit
	Fitness *entity2.Habit
	Work    *entity2.Habit
	Eat     *entity2.Habit
}

func formatDays(ds []*entity2.Day) []dayContainer {
	var res []dayContainer
	for _, d := range ds {
		dc := dayContainer{Date: d.Date}
		for _, h := range d.Habits {
			switch h.Category.Type {
			case entity2.HabitTypeWakeUp:
				dc.Wake = h
			case entity2.HabitTypeFitness:
				dc.Fitness = h
			case entity2.HabitTypeDeepWork:
				dc.Work = h
			case entity2.HabitTypeFood:
				dc.Eat = h
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

func replace(h *entity2.Habit) *entity2.Habit {
	if h == nil {
		return &entity2.Habit{Status: entity2.HabitStatusNoInfo}
	}
	return h
}
