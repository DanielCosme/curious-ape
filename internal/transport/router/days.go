package router

import (
	"github.com/danielcosme/curious-ape/internal/transport/types"
	"github.com/danielcosme/curious-ape/rest"
	"net/http"
)

func (h *Handler) DaysGetAll(rw http.ResponseWriter, r *http.Request) {
	days, err := h.App.DaysGetAll()
	if err != nil {
		rest.ErrBadRequest(rw, r, err)
		return
	}

	daysTransport := []*types.DayTransport{}
	for _, d := range days {
		daysTransport = append(daysTransport, types.DayToTransport(d))
	}
	rest.JSONStatusOk(rw, &envelope{"days": daysTransport})
}
