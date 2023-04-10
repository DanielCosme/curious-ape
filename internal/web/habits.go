package web

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/danielcosme/curious-ape/internal/core/database"
	"github.com/danielcosme/curious-ape/internal/core/entity"
)

func (h *Handler) habitView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		h.notFound(w)
		return
	}

	habit, err := h.App.HabitGetByID(id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			h.notFound(w)
		} else {
			h.serverError(w, err)
		}
		return
	}

	data := h.newTemplateData(r)
	data.Habit = habit
	h.render(w, http.StatusOK, "view.html.tmpl", data)
}

func (h *Handler) habitCreate(w http.ResponseWriter, r *http.Request) {
	d, err := h.App.DayGetByDate(time.Now().AddDate(0, 0, 0))
	if err != nil {
		h.serverError(w, err)
		return
	}

	count := 1
	for count < 5 {
		_, err := h.App.HabitCreate(d, &entity.Habit{
			CategoryID: count,
			Logs: []*entity.HabitLog{
				{
					Success:     (count % 2) == 0,
					Note:        "habit created on loop",
					Origin:      "manual",
					IsAutomated: false,
				},
			},
		})
		if err != nil {
			h.serverError(w, err)
			return
		}
		count++
	}
}
