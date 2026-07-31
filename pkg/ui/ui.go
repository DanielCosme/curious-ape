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
	CLayout        = "layout"
	CNavLink       = "nav-link"
	CNavLinkActive = "nav-link active"
	CBtn           = "btn"
	CBtnNav        = "btn btn-nav"
	CSurface       = "surface"
	CLogEntry      = "log-entry"
	CError         = "error"
	CVersion       = "version"
	CStatusBadge   = "status-badge"
	CSkeleton      = "skeleton"
)

type UIState struct {
	Title           string
	IsAuthenticated bool
	CurrentPath     string
	Version         string
}

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

func StateWithContextUI(ctx context.Context, s *UIState) context.Context {
	return context.WithValue(ctx, CtxState, s)
}

func StateFromContextUI(ctx context.Context) *UIState {
	v, ok := ctx.Value(CtxState).(*UIState)
	if ok {
		return v
	}
	panic("state not set")
}

func StateFromContext(ctx context.Context) *State {
	v, ok := ctx.Value(CtxState).(*State)
	if ok {
		return v
	}
	panic("state not set")
}
