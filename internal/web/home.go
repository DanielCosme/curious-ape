package web

import (
	"net/http"
	"time"

	"github.com/danielcosme/curious-ape/internal/core/entity"
)

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	habits, err := h.App.HabitsGetAll(map[string]string{
		"endDate":   now.Format(entity.ISO8601),
		"startDate": now.AddDate(0, 0, -30).Format(entity.ISO8601),
	})
	if err != nil {
		h.serverError(w, err)
		return
	}

	data := h.newTemplateData(r)
	data.Habits = habits
	h.render(w, http.StatusOK, "home.html.tmpl", data)
}
