package habit

import (
	"net/http"
	"strconv"
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/ui"
	"danicos.dev/daniel/curious-ape/pkg/web"
	"github.com/go-chi/chi/v5"
	// "time"
	// "danicos.dev/daniel/curious-ape/pkg/core"
	// "danicos.dev/daniel/curious-ape/pkg/ui"
)

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) Flip(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		web.ErrInternalServer(err, w)
		return
	}

	habit, err := h.svc.Flip(id)
	if err != nil {
		web.ErrInternalServer(err, w)
		return
	}
	// TODO: I should Send SSE Events
	// 1. For the habit cell.
	// 2. For the re-calculation of the Month Score.
	//
	// OR Just Re-Render the Entire Month back?
	//
	web.Render(w, ui_habitCell(habit))
}

func (h *Handler) Habits(w http.ResponseWriter, r *http.Request) {
	state := ui.StateFromContextUI(r.Context())
	today := core.NewDate(time.Now())

	habitState := &HabitsPageState{}
	for _, month := range today.Months() {
		t := time.Date(today.Time().Year(), month+1, 0, 0, 0, 0, 0, time.UTC)
		d := core.NewDate(t)
		if month == today.Time().Month() {
			d = today
		}

		// Maybe I emit an event here?...

		// TODO: Improve this mess.
		wakeUphabits, err := find(h.svc.db, core.HabitParams{
			From: d.FirstDayOfTheMonth(),
			To:   d.LastDayOfTheMonth(),
			Type: core.HabitTypeWakeUp,
		})
		if err != nil {
			web.ErrInternalServer(err, w)
			return
		}
		fitnessHabits, err := find(h.svc.db, core.HabitParams{
			From: d.FirstDayOfTheMonth(),
			To:   d.LastDayOfTheMonth(),
			Type: core.HabitTypeFitness,
		})
		if err != nil {
			web.ErrInternalServer(err, w)
			return
		}
		workHabits, err := find(h.svc.db, core.HabitParams{
			From: d.FirstDayOfTheMonth(),
			To:   d.LastDayOfTheMonth(),
			Type: core.HabitTypeDeepWork,
		})
		if err != nil {
			web.ErrInternalServer(err, w)
			return
		}
		eatHabits, err := find(h.svc.db, core.HabitParams{
			From: d.FirstDayOfTheMonth(),
			To:   d.LastDayOfTheMonth(),
			Type: core.HabitTypeEatHealthy,
		})
		if err != nil {
			web.ErrInternalServer(err, w)
			return
		}

		m := Month{
			m:       month,
			days:    d.RangeMonth(),
			wakeUp:  wakeUphabits,
			fitness: fitnessHabits,
			work:    workHabits,
			eat:     eatHabits,
		}
		habitState.Months = append(habitState.Months, m)
	}
	web.Render(w, HabitsPage(state, habitState))
}
