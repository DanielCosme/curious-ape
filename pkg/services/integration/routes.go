package integration

import (
	"log/slog"
	"net/http"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/ui"
	"danicos.dev/daniel/curious-ape/pkg/web"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r chi.Router, handler *Handler) error {
	r.Route("/integrations", func(r chi.Router) {
		r.Get("/", handler.IntegrationsPage)
		r.Get("/{name}", handler.IntegrationGet)
	})

	return nil
}

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) IntegrationsPage(w http.ResponseWriter, r *http.Request) {
	state := ui.StateFromContextUI(r.Context())
	var list []core.IntegrationInfo
	for _, integration := range h.svc.sync.IntegrationsList() {
		list = append(list, core.IntegrationInfo{
			Name:   core.ToUpperFist(string(integration)),
			Status: core.IntegrationStatusUnkown,
		})
	}
	web.Render(w, UI_Integrations(state, list))
}

func (h *Handler) IntegrationGet(w http.ResponseWriter, r *http.Request) {
	i, err := h.svc.GetIntegration(chi.URLParam(r, "name"))
	if err != nil {
		web.ErrInternalServer(err, w)
		return
	}
	slog.Info("info", "status", i.Status, "info", i.Info)
	web.Render(w, ui_integration(i))
}
