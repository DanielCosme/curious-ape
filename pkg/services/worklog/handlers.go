package worklog

import (
	"net/http"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/services/day"
	"danicos.dev/daniel/curious-ape/pkg/ui"
	"danicos.dev/daniel/curious-ape/pkg/web"
)

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) worklogPage(w http.ResponseWriter, r *http.Request) {
	days, err := day.Find(h.svc.db, core.DayParams{
		Dates:        web.GetDayParams(r).RangeMonthAll(),
		WithRelation: []core.DayRelations{core.DayRelationDeepWorkLogs},
	})
	if err != nil {
		web.ErrInternalServer(err, w)
		return
	}

	state := ui.StateFromContext(r.Context())
	web.Render(w, UI_WorkLog(state, days))
}
