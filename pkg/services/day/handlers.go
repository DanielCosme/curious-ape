package day

import (
	"errors"
	"log/slog"
	"net/http"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/event"
	"danicos.dev/daniel/curious-ape/pkg/ui"
	"danicos.dev/daniel/curious-ape/pkg/web"
	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
	"github.com/starfederation/datastar-go/datastar"
)

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) index(w http.ResponseWriter, r *http.Request) {
	d := web.GetDayParams(r)
	days, err := h.svc.Month(d, core.DESC)
	if err != nil {
		web.ErrInternalServer(err, w)
		return
	}

	s := ui.StateFromContext(r.Context())
	web.Render(w, UI_Index(s, days))
}

func (h *Handler) streamSSE(w http.ResponseWriter, r *http.Request) {
	slog.Info("Day: SSE stream open")
	sse := datastar.NewSSE(w, r)

	ch := make(chan *nats.Msg)
	subs, err := h.svc.nats.ChanSubscribe(event.DaySynced, ch)
	if err != nil {
		web.ErrInternalServer(err, w)
		return
	}
	defer subs.Unsubscribe()

	for {
		select {
		case <-r.Context().Done():
			slog.Warn("Day: SSE stream closed")
			return
		case msg := <-ch:
			day, err := h.svc.GetOrCreate(core.DateDecode(msg.Data))
			if err != nil {
				web.ErrInternalServer(err, w)
				return
			}

			err = sse.PatchElementGostar(UI_day(day))
			if err != nil {
				web.ErrInternalServer(err, w)
				return
			}
		}
	}
}

func (h *Handler) sync(w http.ResponseWriter, r *http.Request) {
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
