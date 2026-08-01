package day

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/event"
	"danicos.dev/daniel/curious-ape/pkg/ui"
	"danicos.dev/daniel/curious-ape/pkg/web"
	"github.com/go-chi/chi/v5"
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

	s := ui.StateFromContextUI(r.Context())
	web.Render(w, UI_Index(s, days))
}

func (h *Handler) streamSSE(w http.ResponseWriter, r *http.Request) {
	slog.Info("stream endpoint created")
	sse := datastar.NewSSE(w, r)

	subs, err := h.svc.nats.SubscribeSync(event.DaySynced)
	if err != nil {
		web.ErrInternalServer(err, w)
		return
	}

	for {
		msg, err := subs.NextMsg(time.Hour * 12)
		if err != nil {
			web.ErrInternalServer(err, w)
			return
		}

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
		continue
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
