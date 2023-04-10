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
	wc.Server.Handler = wc.Routes()

	wc.App.Log.InfoP("HTTP server listening", log.Prop{"addr": wc.Server.Addr})
	return wc.Server.ListenAndServe()
}

func (wc *WebClient) Routes() http.Handler {
	return ChiRoutes(wc.App)
}
