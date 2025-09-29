package api

import (
	"github.com/danielcosme/curious-ape/pkg/application"
)

type API struct {
	App     *application.App
	Version string
}

func NewApi(app *application.App, version string) *API {
	return &API{
		App:     app,
		Version: version,
	}
}
