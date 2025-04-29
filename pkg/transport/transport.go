package transport

import (
	"github.com/danielcosme/curious-ape/pkg/application"
	"net/http"
)

type Transport struct {
	app     *application.Application
	version string
}

func New(a *application.Application, v string) *Transport {
	return &Transport{
		app:     a,
		version: v,
	}
}

func (t *Transport) HandlerInfo(w http.ResponseWriter, r *http.Request) {
	info := map[string]string{
		"version":     t.version,
		"environment": t.app.Env,
	}
	JSONOK(w, envelope{"info": info}, nil)
}
