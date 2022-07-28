package routes

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/rest"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) SyncByDate(rw http.ResponseWriter, r *http.Request) {
	if param := chi.URLParam(r, "resourceToSync"); isValidResource(param) {
		startDate, err := entity.ParseDate(chi.URLParam(r, "start"))
		if err != nil {
			rest.ErrBadRequest(rw, err.Error())
			return
		}

		switch param {
		case "sleep":
			err = h.App.SyncFitbitSleepLog(startDate)
		case "fitness":
			err = h.App.SyncFitnessLog(startDate)
		}
		JsonCheckError(rw, r, http.StatusOK, envelope{"success": "ok"}, err)
		return
	}

	rest.ErrBadRequest(rw, "no valid resource to sync")
}

func (h *Handler) SyncByDateRange(rw http.ResponseWriter, r *http.Request) {
	if param := chi.URLParam(r, "resourceToSync"); isValidResource(param) {
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

		switch param {
		case "sleep":
			err = h.App.SyncSleepByDateRange(startDate, endDate)
		case "fitness":
			err = h.App.SyncFitnessByDateRAnge(startDate, endDate)
		}
		JsonCheckError(rw, r, http.StatusOK, envelope{"success": "ok"}, err)
		return
	}

	rest.ErrBadRequest(rw, "no valid resource to sync")
}

func (h *Handler) Sync(rw http.ResponseWriter, r *http.Request) {
	var err error
	if param := chi.URLParam(r, "resourceToSync"); isValidResource(param) {
		switch param {
		case "sleep":
			err = h.App.SyncFitbitSleep()
		case "fitness":
			// err = h.App.SyncFitnessLog()
		}
		JsonCheckError(rw, r, http.StatusOK, envelope{"success": "ok"}, err)
		return
	}

	rest.ErrBadRequest(rw, "no valid resource to sync")
}

func isValidResource(s string) bool {
	switch s {
	case "sleep", "fitness":
		return true
	}
	return false
}
