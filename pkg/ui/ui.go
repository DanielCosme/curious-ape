package ui

import (
	// "github.com/eduardolat/gomponents-lucide"
	"context"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"

	"danicos.dev/daniel/curious-ape/pkg/core"
)

type CtxKey string

const CtxState CtxKey = "ui_state"

// Classes (central place for reusable class names)
const (
	cLayout        = "layout"
	cNavLink       = "nav-link"
	cNavLinkActive = "nav-link active"
	CBtn           = "btn"
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
	Integrations  []core.IntegrationInfo
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

func StateWithContext(ctx context.Context, s *State) context.Context {
	return context.WithValue(ctx, CtxState, s)
}

func StateFromContext(ctx context.Context) *State {
	v, ok := ctx.Value(CtxState).(*State)
	if ok {
		return v
	}
	panic("state not set")
}
