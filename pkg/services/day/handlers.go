package day

import (
	"log/slog"
	"net/http"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/ui"
	"danicos.dev/daniel/curious-ape/pkg/web"
)

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	paramsDate := web.GetDayParams(r)
	slog.Info(paramsDate.String())
	days, err := h.svc.Month(paramsDate, core.DESC)
	if err != nil {
		web.ErrInternalServer(err, w)
		return
	}

	s := ui.StateFromContextUI(r.Context())
	web.Render(w, UI_Index(s, days))
}
