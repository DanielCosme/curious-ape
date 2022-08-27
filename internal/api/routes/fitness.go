package routes

import (
	"github.com/danielcosme/curious-ape/internal/api/types"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/rest"
	"net/http"
)

func (h *Handler) FitnessGetAll(rw http.ResponseWriter, r *http.Request) {
	fls, err := h.App.FitnessFindLogs(entity.FitnessLogFilter{})
	JsonCheckError(rw, http.StatusOK, envelope{"fitness_logs": types.FromFitnessLogToTransportSlice(fls)}, err)
}

func (h *Handler) FitnessGetForDate(rw http.ResponseWriter, r *http.Request) {
	day := r.Context().Value("day").(*entity.Day)
	fls, err := h.App.FitnessFindLogs(entity.FitnessLogFilter{DayID: []int{day.ID}})
	JsonCheckError(rw, http.StatusOK, envelope{"fitness_logs": types.FromFitnessLogToTransportSlice(fls)}, err)
}

func (h *Handler) FitnessDelete(rw http.ResponseWriter, r *http.Request) {
	fitnessLog := r.Context().Value("fitnessLog").(*entity.FitnessLog)
	err := h.App.FitnessDeleteLog(fitnessLog)
	JsonCheckError(rw, http.StatusOK, envelopeSuccess(), err)
}

func (h *Handler) FitnessGet(rw http.ResponseWriter, r *http.Request) {
	fitnessLog := r.Context().Value("fitnessLog").(*entity.FitnessLog)
	rest.JSON(rw, http.StatusOK, envelope{"fitness_logs": types.FromFitnessLogToTransport(fitnessLog)})
}

func (h *Handler) FitnessCreate(rw http.ResponseWriter, r *http.Request) {
	day := r.Context().Value("day").(*entity.Day)

	var data *types.FitnessLogTransport
	if err := rest.ReadJSON(r, &data); err != nil {
		h.ErrInternalServerError(rw, err)
		return
	}
	fitnessLog, err := data.ToFitnessLog(day)
	if err != nil {
		h.ErrInternalServerError(rw, err)
		return
	}

	fitnessLog, err = h.App.FitnessCreateLogFromApi(fitnessLog)
	JsonCheckError(rw, http.StatusCreated, envelope{"fitness_logs": types.FromFitnessLogToTransport(fitnessLog)}, err)
}

func (h *Handler) FitnessUpdate(rw http.ResponseWriter, r *http.Request) {
	fitnessLog := r.Context().Value("fitnessLog").(*entity.FitnessLog)

	var payload *types.FitnessLogTransport
	err := rest.ReadJSON(r, &payload)
	if err != nil {
		rest.ErrInternalServer(rw)
		return
	}
	data, err := payload.ToFitnessLog(fitnessLog.Day)
	if err != nil {
		rest.ErrInternalServer(rw)
		return
	}

	fitnessLog, err = h.App.UpdateFitnessLog(fitnessLog, data)
	JsonCheckError(rw, http.StatusOK, envelope{"fitness_logs": types.FromFitnessLogToTransport(fitnessLog)}, err)
}

func (h *Handler) ErrInternalServerError(rw http.ResponseWriter, err error) {
	h.App.Log.Error(err)
	rest.ErrInternalServer(rw)
}
