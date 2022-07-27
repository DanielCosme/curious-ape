package routes

import (
	"net/http"
)

func (h *Handler) FitnessGet(rw http.ResponseWriter, r *http.Request) {
	h.App.Log.Debug("we are in the fitness handler")
	err := h.App.FitnessGet()
	JsonCheckError(rw, r, http.StatusOK, envelope{"success": "ok"}, err)
}
