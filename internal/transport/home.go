package transport

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"net/http"
	"time"
)

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	ds, err := h.App.DaysMonth()
	if err != nil {
		h.serverError(w, err)
		return
	}
	data := h.newTemplateData(r)
	data.Days = formatDays(ds)
	h.render(w, http.StatusOK, "home.gohtml", data)
}

type dayContainer struct {
	Date    time.Time
	Wake    *entity.Habit
	Fitness *entity.Habit
	Work    *entity.Habit
	Eat     *entity.Habit
}

func formatDays(ds []*entity.Day) []dayContainer {
	var res []dayContainer
	for _, d := range ds {
		dc := dayContainer{Date: d.Date}
		for _, h := range d.Habits {
			switch h.Category.Type {
			case entity.HabitTypeWakeUp:
				dc.Wake = h
			case entity.HabitTypeFitness:
				dc.Fitness = h
			case entity.HabitTypeDeepWork:
				dc.Work = h
			case entity.HabitTypeFood:
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

func replace(h *entity.Habit) *entity.Habit {
	if h == nil {
		return &entity.Habit{Status: entity.HabitStatusNoInfo}
	}
	return h
}
