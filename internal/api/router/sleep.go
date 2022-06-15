package router

import (
	"github.com/danielcosme/curious-ape/internal/api/types"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/rest"
	"net/http"
)

func (h *Handler) SleepLogs(rw http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case http.MethodGet:
		day := r.Context().Value("day").(*entity.Day)
		sls := []*entity.SleepLog{}
		if day != nil { // we want a single day
			sls, err = h.App.GetSleepLogsForDay(day)
		} else { // we just get them all
			sls, err = h.App.GetAllSleepLogs()
		}
		JsonCheckError(rw, r, http.StatusOK, envelope{"sleep_logs": types.FromSleepLogToTransportSlice(sls)}, err)
	case http.MethodPost:
		switch r.Header.Get("X-APE-ACTION") {
		case "sync":
			start := r.Header.Get("X-APE-DATE")
			if start == "" {
				rest.ErrBadRequest(rw, "missing start date")
				return
			}
			startDate, err := entity.ParseDate(start)
			if err != nil {
				rest.ErrBadRequest(rw, err.Error())
				return
			}

			end := r.Header.Get("X-APE-DATE-END")
			if end != "" { // sync from start to end
				endDate, err := entity.ParseDate(end)
				if err != nil {
					rest.ErrBadRequest(rw, err.Error())
					return
				}

				err = h.App.SyncSleepLogs(startDate, endDate)
			} else { // sync only that day
				err = h.App.SyncSleepLog(startDate)
			}
			JsonCheckError(rw, r, http.StatusOK, envelope{"success": "ok"}, err)
		default:
			rest.ErrBadRequest(rw, envelope{"message": "action not specified in header", "missing header": "X-APE-ACTION"})
		}
	default:
		rest.ErrNotAllowed(rw)
	}
}
