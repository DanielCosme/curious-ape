package day

import (
	"errors"
	"net/http"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/event"
	"danicos.dev/daniel/curious-ape/pkg/ui"
	"danicos.dev/daniel/curious-ape/pkg/web"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	days, err := h.svc.Month(core.NewDateToday(), core.DESC)
	if err != nil {
		web.ErrInternalServer(err, w)
		return
	}

	s := ui.StateFromContextUI(r.Context())
	web.Render(w, UI_Index(s, days))
}

func (h *Handler) Sync(w http.ResponseWriter, r *http.Request) {
	// NOET: if this grows, (date param in URL) create a middleware that sets the date.
	dateParam := chi.URLParam(r, "date")
	if dateParam == "" {
		web.ErrBadRequest(errors.New("date param is empty"), w)
		return
	}
	date, err := core.NewDateFromISO8601(dateParam)
	if err != nil {
		web.ErrBadRequest(err, w)
		return
	}

	h.svc.nats.Publish(event.DaySync, date.Enc())
	w.WriteHeader(http.StatusAccepted)
}
