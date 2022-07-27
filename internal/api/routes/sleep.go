package routes

import (
	"github.com/danielcosme/curious-ape/internal/api/types"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"net/http"
)

func (h *Handler) SleepGetByDate(rw http.ResponseWriter, r *http.Request) {
	day := r.Context().Value("day").(*entity.Day)
	sls, err := h.App.GetSleepLogsForDay(day)
	JsonCheckError(rw, r, http.StatusOK, envelope{"sleep_logs": types.FromSleepLogToTransportSlice(sls)}, err)
}

func (h *Handler) SleepGetAll(rw http.ResponseWriter, r *http.Request) {
	sls, err := h.App.GetAllSleepLogs()
	JsonCheckError(rw, r, http.StatusOK, envelope{"sleep_logs": types.FromSleepLogToTransportSlice(sls)}, err)
}
