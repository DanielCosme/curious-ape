package web

import (
	"net/http"

	"github.com/danielcosme/curious-ape/internal/core/entity"
)

func (h *Handler) habitView(w http.ResponseWriter, r *http.Request) {
	habit := r.Context().Value("habit").(*entity.Habit)

	data := h.newTemplateData(r)
	data.Habit = habit
	h.render(w, http.StatusOK, "view.html.tmpl", data)
}

func (h *Handler) habitCreateForm(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) habitCreate(w http.ResponseWriter, r *http.Request) {
}
