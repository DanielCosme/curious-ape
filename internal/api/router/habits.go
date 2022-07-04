package router

import (
	"github.com/danielcosme/curious-ape/internal/api/types"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/rest"
	"net/http"
	"time"
)

func (h *Handler) HabitsGetCategories(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		categories, err := h.App.HabitsGetCategories()
		JsonCheckError(rw, r, http.StatusOK, &envelope{"categories": categories}, err)
	}
}

func (h *Handler) HabitGet(rw http.ResponseWriter, r *http.Request) {
	habit := r.Context().Value("habit").(*entity.Habit)
	rest.JSONStatusOk(rw, envelope{"habit": types.FromHabitToTransport(habit)})
}

func (h *Handler) HabitCreate(rw http.ResponseWriter, r *http.Request) {
	day := r.Context().Value("day").(*entity.Day)

	var data *types.HabitTransport
	err := rest.ReadJSON(r, &data)
	if err != nil {
		rest.ErrInternalServer(rw)
		return
	}

	newHabit, err := h.App.HabitCreate(day, data.ToHabit())
	JsonCheckError(rw, r, http.StatusCreated, envelope{"habit": types.FromHabitToTransport(newHabit)}, err)
}

func (h *Handler) HabitUpdate(rw http.ResponseWriter, r *http.Request) {
	habit := r.Context().Value("habit").(*entity.Habit)

	var data *types.HabitTransport
	err := rest.ReadJSON(r, &data)
	if err != nil {
		rest.ErrInternalServer(rw)
		return
	}

	habitUpdated, err := h.App.HabitFullUpdate(habit, data.ToHabit())
	JsonCheckError(rw, r, http.StatusOK, envelope{"habit": types.FromHabitToTransport(habitUpdated)}, err)
}

func (h *Handler) HabitDelete(rw http.ResponseWriter, r *http.Request) {
	habit := r.Context().Value("habit").(*entity.Habit)
	err := h.App.HabitDelete(habit)
	JsonCheckError(rw, r, http.StatusOK, nil, err)
}

func (h *Handler) HabitsGetAll(rw http.ResponseWriter, r *http.Request) {
	// should I send a filter here?

	// I could have a day object
	// Create with URL param  ?day=2022-01-02

	// Get all
	// Get all by day
	// Get all from day to day
	// day := r.Context().Value("day").(*entity.Day)
	// entity.HabitFilter{
	// 	ID:         nil,
	// 	DayID:      []int{day.ID},
	// 	CategoryID: nil,
	// }

	hs, err := h.App.HabitsGetAll(time.Now(), time.Now())
	if err != nil {
		rest.ErrNotFound(rw)
		return
	}

	habits := []*types.HabitTransport{}
	for _, habit := range hs {
		habits = append(habits, types.FromHabitToTransport(habit))
	}

	rest.JSONStatusOk(rw, envelope{"habits": hs})
	return
}
