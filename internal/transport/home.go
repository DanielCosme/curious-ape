package transport

import (
	"net/http"
	"sort"
	"time"

	"github.com/danielcosme/curious-ape/internal/application"
	"github.com/labstack/echo/v4"

	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/view"
)

func (t *Transport) home(c echo.Context) error {
	var d core.Date
	var err error
	if dayParam := c.QueryParam("day"); dayParam == "" {
		d = core.NewDate(time.Now())
	} else {
		d, err = core.DateFromISO8601(dayParam)
		if err != nil {
			return errClientError(err)
		}
	}
	ds, err := t.App.DaysMonth(d)
	if err != nil {
		return errServer(err)
	}
	td := t.newTemplateData(c.Request())
	sort.Sort(application.DaysSlice(ds))
	days := formatDays(ds)
	return t.RenderTempl(http.StatusOK, c, view.Home(td, days))
}

func formatDays(ds []core.Day) []view.DaySummary {
	var res []view.DaySummary
	for _, d := range ds {
		dc := view.DaySummary{Date: d.Date.Time()}
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
			if h.State() == core.HabitStateDone {
				dc.Score++
			}
		}
		dc.Wake = replace(dc.Wake)
		dc.Wake.Category.Type = core.HabitTypeWakeUp
		dc.Fitness = replace(dc.Fitness)
		dc.Fitness.Category.Type = core.HabitTypeExercise
		dc.Work = replace(dc.Work)
		dc.Work.Category.Type = core.HabitTypeDeepWork
		dc.Eat = replace(dc.Eat)
		dc.Eat.Category.Type = core.HabitTypeEatHealthy
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
