package deadline

import (
	"log/slog"
	"net/http"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/ui"
	"danicos.dev/daniel/curious-ape/pkg/web"
)

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) deadlinePage(w http.ResponseWriter, r *http.Request) {
	ds, err := h.svc.List()
	if err != nil {
		web.ErrInternalServer(err, w)
		return
	}
	deadlinesState := ui.DeadlineState{DS: ds}
	state := ui.StateFromContext(r.Context())
	web.Render(w, UI_Deadlines(state, deadlinesState))
}

func (h *Handler) newDeadlinePage(w http.ResponseWriter, r *http.Request) {
	state := ui.StateFromContext(r.Context())
	web.Render(w, DeadlineForm(state, ui.DeadlineState{}))
}

func (h *Handler) newDeadlinePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	state := ui.StateFromContext(r.Context())

	var recurring bool
	if r.PostForm.Get("recurrent") == "on" {
		recurring = true
	}
	date, err := core.NewDateFromISO8601(r.PostForm.Get("end_date"))
	if err == nil {
		_, err := h.svc.NewDeadline(core.Deadline{
			Title:     r.PostForm.Get("title"),
			StartDate: core.NewDateToday(),
			EndDate:   date,
			Recurring: recurring,
		})
		if err == nil {
			web.Redirect(w, "/deadlines")
			return
		}
		deadlineState := ui.DeadlineState{Err: err}
		web.Render(w, DeadlineForm(state, deadlineState))
		return
	}

	web.ErrBadRequest(err, w)
}

func (svc *Service) NewDeadline(params core.Deadline) (core.Deadline, error) {
	res, err := create(svc.db, params)
	if err != nil {
		return params, err
	}
	slog.Info("Deadline created",
		"Title", res.Title,
		"End Date", res.EndDate.String(),
		"recurring", res.Recurring,
	)
	return res, nil
}
