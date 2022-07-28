package routes

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/rest"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) SyncSleepByDate(rw http.ResponseWriter, r *http.Request) {
	startDate, err := entity.ParseDate(chi.URLParam(r, "start"))
	if err != nil {
		rest.ErrBadRequest(rw, err.Error())
		return
	}

	err = h.App.SyncFitbitSleepLog(startDate)
	JsonCheckError(rw, r, http.StatusOK, envelope{"success": "ok"}, err)
}

func (h *Handler) SyncSleepByDateRange(rw http.ResponseWriter, r *http.Request) {
	startDate, err := entity.ParseDate(chi.URLParam(r, "startDate"))
	if err != nil {
		rest.ErrBadRequest(rw, err.Error())
		return
	}
	endDate, err := entity.ParseDate(chi.URLParam(r, "endDate"))
	if err != nil {
		rest.ErrBadRequest(rw, err.Error())
		return
	}

	err = h.App.SyncFitbitSleepByDateRange(startDate, endDate)
	JsonCheckError(rw, r, http.StatusOK, envelope{"success": "ok"}, err)
}

func (h *Handler) SyncSleep(rw http.ResponseWriter, r *http.Request) {
	err := h.App.SyncFitbitSleep()
	JsonCheckError(rw, r, http.StatusOK, envelope{"success": "ok"}, err)
}
