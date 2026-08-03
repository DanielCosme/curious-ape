package ui

import (
	"context"

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

type State struct {
	Title           string
	IsAuthenticated bool
	CurrentPath     string
	Version         string
}

type DeadlineState struct {
	Err error
	DS  []core.Deadline
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
