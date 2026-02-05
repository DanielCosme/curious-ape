package ui

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"

	"github.com/danielcosme/curious-ape/pkg/application"
	"github.com/danielcosme/curious-ape/pkg/core"
)

type State struct {
	Version       string
	Authenticated bool
	Days          []core.Day
	Integrations  []application.IntegrationInfo
}

func a(path, name string) Node {
	return A(Href(path), Text(name))
}

func blockDisplay() Node {
	return Style("display: block;")
}
