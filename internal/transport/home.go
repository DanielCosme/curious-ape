package transport

import (
	"github.com/danielcosme/curious-ape/internal/application"
	"github.com/labstack/echo/v4"
	"net/http"
	"sort"
	"time"

	"github.com/danielcosme/curious-ape/internal/core"
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
	data := t.newTemplateData(c.Request())
	sort.Sort(application.DaysSlice(ds))
	data.Days = formatDays(ds)
	c.Set("page", pageHome)
	return c.Render(http.StatusOK, pageHome, data)
}

type dayContainer struct {
	Date    time.Time
	Wake    core.Habit
	Fitness core.Habit
	Work    core.Habit
	Eat     core.Habit
	Score   int
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
