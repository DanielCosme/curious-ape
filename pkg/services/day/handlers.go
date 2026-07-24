package day

import (
	"net/http"

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
	s := ui.StateFromContextUI(r.Context())
	web.Render(w, UI_Index(s, nil))
}
