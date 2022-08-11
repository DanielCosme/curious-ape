package routes

import (
	"net/http"
)

func (h *Handler) TogglGetProjects(rw http.ResponseWriter, r *http.Request) {
	ps, err := h.App.TogglGetProjects()
	JsonCheckError(rw, r, http.StatusOK, envelope{"toggl_projects": ps}, err)
}

func (h *Handler) TogglAssignProjectsToGoal(rw http.ResponseWriter, r *http.Request) {
	err := h.App.TogglAssignProjectsToGoal(r.Form.Get("ids"))
	JsonCheckError(rw, r, http.StatusOK, envelope{"success": "ok"}, err)
}
