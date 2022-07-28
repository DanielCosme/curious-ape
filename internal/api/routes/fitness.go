package routes

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"net/http"
	"time"
)

func (h *Handler) FitnessGet(rw http.ResponseWriter, r *http.Request) {
	h.App.Log.Debug("we are in the fitness handler")
	err := h.App.SyncFitnessLog(time.Now())
	JsonCheckError(rw, r, http.StatusOK, envelope{"success": "ok"}, err)
}

func (h *Handler) FitnessGetByDate(rw http.ResponseWriter, r *http.Request) {
	day := r.Context().Value("day").(*entity.Day)
	fls, err := h.App.GetFitnessLogsForDay(day)
	// TODO implement the transport functionality
	JsonCheckError(rw, r, http.StatusCreated, fls, err)
}
