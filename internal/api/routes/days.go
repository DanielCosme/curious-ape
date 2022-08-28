package routes

import (
	"github.com/danielcosme/curious-ape/internal/api/types"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/rest"
	"net/http"
)

func (h *Handler) DaysGetAll(rw http.ResponseWriter, r *http.Request) {
	days, err := h.App.DaysGetAll()
	if err != nil {
		rest.ErrBadRequest(rw, err)
		return
	}

	daysTransport := []*types.DayTransport{}
	for _, d := range days {
		daysTransport = append(daysTransport, types.DayToTransport(d))
	}
	rest.JSONStatusOk(rw, &envelope{"days": daysTransport})
}

func (h *Handler) DayGetByDate(rw http.ResponseWriter, r *http.Request) {
	day := r.Context().Value("day").(*entity.Day)
	rest.JSONStatusOk(rw, envelope{"day": types.DayToTransport(day)})
}

func (h *Handler) DayUpdate(rw http.ResponseWriter, r *http.Request) {
	day := r.Context().Value("day").(*entity.Day)

	var data *types.DayTransport
	err := rest.ReadJSON(r, &data)
	if err != nil {
		h.ErrInternalServerError(rw, err)
		return
	}

	day, err = h.App.DayUpdate(day, data.ToDay())
	JsonCheckError(rw, http.StatusOK, envelope{"day": types.DayToTransport(day)}, err)
}
