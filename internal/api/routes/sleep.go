package routes

import (
	"github.com/danielcosme/curious-ape/internal/api/types"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/rest"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *Handler) SleepGetByDate(rw http.ResponseWriter, r *http.Request) {
	day := r.Context().Value("day").(*entity.Day)
	sls, err := h.App.GetSleepLogsForDay(day)
	JsonCheckError(rw, r, http.StatusOK, envelope{"sleep_logs": types.FromSleepLogToTransportSlice(sls)}, err)
}

func (h *Handler) SleepGetAll(rw http.ResponseWriter, r *http.Request) {
	sls, err := h.App.GetSleepLogs()
	JsonCheckError(rw, r, http.StatusOK, envelope{"sleep_logs": types.FromSleepLogToTransportSlice(sls)}, err)
}

func (h *Handler) SleepDeleteByID(rw http.ResponseWriter, r *http.Request) {
	var id int
	var err error
	if idStr := chi.URLParam(r, "id"); idStr != "" {
		id, err = strconv.Atoi(idStr)
		if err != nil {
			rest.ErrInternalServer(rw)
			return
		}
	}
	err = h.App.SleepLogDeleteByID(id)
	JsonCheckError(rw, r, http.StatusOK, envelope{"success": "ok"}, err)
}

func (h *Handler) SleepCreate(rw http.ResponseWriter, r *http.Request) {
	day := r.Context().Value("day").(*entity.Day)

	var data *types.SleepLogTransport
	err := rest.ReadJSON(r, &data)
	if err != nil {
		rest.ErrInternalServer(rw)
		return
	}
	sleepLog, err := data.ToSleepLog(day)
	if err != nil {
		rest.ErrInternalServer(rw)
		return
	}

	sleepLog, err = h.App.SleepFromRestCreate(sleepLog)
	JsonCheckError(rw, r, http.StatusOK, envelope{"sleep_logs": types.FromSleepLogToTransport(sleepLog)}, err)
}
