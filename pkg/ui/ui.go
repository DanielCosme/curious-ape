package ui

import (
	// "github.com/eduardolat/gomponents-lucide"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"

	"danicos.dev/daniel/curious-ape/pkg/application"
	"danicos.dev/daniel/curious-ape/pkg/core"
)

// Classes (central place for reusable class names)
const (
	cLayout        = "layout"
	cNavLink       = "nav-link"
	cNavLinkActive = "nav-link active"
	cBtn           = "btn"
	cBtnNav        = "btn btn-nav"
	cSurface       = "surface"
	cLogEntry      = "log-entry"
	cError         = "error"
	cVersion       = "version"
	cStatusBadge   = "status-badge"
	cSkeleton      = "skeleton"
)

type State struct {
	Version       string
	Authenticated bool
	CurrentPath   string
	DaysYear      [][]core.Day
	Days          []core.Day
	Integrations  []application.IntegrationInfo
	Deadlines     DeadlineState
}

type DeadlineState struct {
	Err error
	DS  []core.Deadline
}

func a(path, name string) Node {
	return A(Href(path), Text(name))
}

func blockDisplay() Node {
	return Style("display: block;")
}
