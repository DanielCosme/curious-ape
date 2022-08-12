package routes

import (
	"github.com/danielcosme/curious-ape/internal/api/types"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/rest"
	"net/http"
)

func (h *Handler) SleepGetForDate(rw http.ResponseWriter, r *http.Request) {
	day := r.Context().Value("day").(*entity.Day)
	sls, err := h.App.GetSleepLogs(entity.SleepLogFilter{DayID: []int{day.ID}})
	JsonCheckError(rw, http.StatusOK, envelope{"sleep_logs": types.FromSleepLogToTransportSlice(sls)}, err)
}

func (h *Handler) SleepGet(rw http.ResponseWriter, r *http.Request) {
	sleepLog := r.Context().Value("sleepLog").(*entity.SleepLog)
	rest.JSON(rw, http.StatusOK, envelope{"sleep_logs": types.FromSleepLogToTransport(sleepLog)})
}

func (h *Handler) SleepGetAll(rw http.ResponseWriter, r *http.Request) {
	sls, err := h.App.GetSleepLogs(entity.SleepLogFilter{})
	JsonCheckError(rw, http.StatusOK, envelope{"sleep_logs": types.FromSleepLogToTransportSlice(sls)}, err)
}

func (h *Handler) SleepDelete(rw http.ResponseWriter, r *http.Request) {
	sleepLog := r.Context().Value("sleepLog").(*entity.SleepLog)
	err := h.App.DeleteSleepByID(sleepLog.ID)
	JsonCheckError(rw, http.StatusOK, envelopeSuccess(), err)
}

func (h *Handler) SleepUpdate(rw http.ResponseWriter, r *http.Request) {
	sleepLog := r.Context().Value("sleepLog").(*entity.SleepLog)

	var payload *types.SleepLogTransport
	err := rest.ReadJSON(r, &payload)
	if err != nil {
		rest.ErrInternalServer(rw)
		return
	}
	data, err := payload.ToSleepLog(sleepLog.Day)
	if err != nil {
		rest.ErrInternalServer(rw)
		return
	}

	sleepLog, err = h.App.UpdateSleep(sleepLog, data)
	JsonCheckError(rw, http.StatusOK, envelope{"sleep_logs": types.FromSleepLogToTransport(sleepLog)}, err)
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

	sleepLog, err = h.App.CreateSleepFromApi(sleepLog)
	JsonCheckError(rw, http.StatusCreated, envelope{"sleep_logs": types.FromSleepLogToTransport(sleepLog)}, err)
}
