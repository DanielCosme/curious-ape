package web

import (
	"net/http"

	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/go-sdk/log"
)

type WebClient struct {
	App    *application.App
	Server *http.Server
}

func (wc *WebClient) ListenAndServe() error {
	h := &Handler{App: wc.App}
	tc, err := newTemplateCache()
	if err != nil {
		return err
	}
	h.templateCache = tc

	wc.Server.Handler = wc.Routes(h)

	wc.App.Log.InfoP("HTTP server listening", log.Prop{"addr": wc.Server.Addr})
	return wc.Server.ListenAndServe()
}

func (wc *WebClient) Routes(h *Handler) http.Handler {
	return ChiRoutes(h)
}
