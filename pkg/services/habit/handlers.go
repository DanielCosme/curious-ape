package habit

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/ui"
	"danicos.dev/daniel/curious-ape/pkg/web"
	"github.com/go-chi/chi/v5"
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

	month := habit.Date.Time.Month()
	habitsMonth, err := find(h.svc.db, core.HabitParams{
		From: habit.Date.FirstDayOfTheMonth(),
		To:   habit.Date.LastDayOfTheMonth(),
	})
	if err != nil {
		web.ErrInternalServer(err, w)
		return
	}

	m := toUIMonth(month, habitsMonth)
	web.Render(w, HabitsGrid(m))
}

func (h *Handler) HabitsPage(w http.ResponseWriter, r *http.Request) {
	state := ui.StateFromContext(r.Context())
	today := core.NewDate(time.Now())

	habitState := &HabitsPageState{}
	for _, month := range today.MonthsOfYear() {
		t := time.Date(today.Time.Year(), month+1, 0, 0, 0, 0, 0, time.UTC)
		d := core.NewDate(t)
		if month == today.Time.Month() {
			d = today
		}

		habitsMonth, err := find(h.svc.db, core.HabitParams{
			From: d.FirstDayOfTheMonth(),
			To:   d.LastDayOfTheMonth(),
		})
		if err != nil {
			web.ErrInternalServer(err, w)
			return
		}
		if len(habitsMonth) == 0 {
			continue
		}

		m := toUIMonth(month, habitsMonth)
		habitState.Months = append(habitState.Months, m)
	}

	web.Render(w, HabitsPage(state, habitState))
}

func toUIMonth(month time.Month, hs []core.Habit) Month {
	m := Month{m: month}
	for _, h := range hs {
		switch h.Type {
		case core.HabitTypeWakeUp:
			m.days = append(m.days, h.Date)
			m.wakeUp = append(m.wakeUp, h)
		case core.HabitTypeFitness:
			m.fitness = append(m.fitness, h)
		case core.HabitTypeDeepWork:
			m.work = append(m.work, h)
		case core.HabitTypeEatHealthy:
			m.eat = append(m.eat, h)
		default:
			panic(fmt.Sprintf("unexpected core.HabitType: %#v", h.Type))
		}

		if h.State == core.HabitStateDone {
			m.score++
		}
	}

	return m
}
