package router

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/transport/types"
	"github.com/danielcosme/curious-ape/rest"
	"net/http"
)

func (h *Handler) SleepLogs(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		day := r.Context().Value("day").(*entity.Day)
		if day != nil { // we want a single day
			sls, err := h.App.GetSleepLogsForDay(day)
			JsonCheckError(rw, r, http.StatusCreated, envelope{"sleep_logs": types.FromSleepLogToTransportSlice(sls)}, err)
		} else { // we assume a get all

		}
	default:
	rest.ErrNotAllowed(rw)
	}
}
