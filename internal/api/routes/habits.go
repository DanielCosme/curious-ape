package routes

import (
	"net/http"

	"github.com/danielcosme/curious-ape/internal/api/types"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/rest"
)

func (h *Handler) HabitsGetCategories(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		categories, err := h.App.HabitsGetCategories()
		JsonCheckError(rw, http.StatusOK, &envelope{"categories": categories}, err)
	}
}

func (h *Handler) HabitGet(rw http.ResponseWriter, r *http.Request) {
	habit := r.Context().Value("habit").(*entity.Habit)
	rest.JSONStatusOk(rw, envelope{"habit": types.FromHabitToTransport(habit)})
}

func (h *Handler) HabitsGetByDay(rw http.ResponseWriter, r *http.Request) {
	day := r.Context().Value("day").(*entity.Day)
	habits, err := h.App.HabitsGetByDay(day)
	JsonCheckError(rw, http.StatusOK, &envelope{"habits": types.FromHabitToTransportSlice(habits)}, err)
}

func (h *Handler) HabitCreate(rw http.ResponseWriter, r *http.Request) {
	// day := r.Context().Value("day").(*entity.Day)

	var data *types.HabitTransport
	err := rest.ReadJSON(r, &data)
	if err != nil {
		rest.ErrInternalServer(rw)
		return
	}

	newHabit, err := h.App.HabitCreate(nil)
	JsonCheckError(rw, http.StatusCreated, envelope{"habit": types.FromHabitToTransport(newHabit)}, err)
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
	JsonCheckError(rw, http.StatusOK, envelope{"habit": types.FromHabitToTransport(habitUpdated)}, err)
}

func (h *Handler) HabitDelete(rw http.ResponseWriter, r *http.Request) {
	habit := r.Context().Value("habit").(*entity.Habit)
	err := h.App.HabitDelete(habit)
	JsonCheckError(rw, http.StatusOK, nil, err)
}

func (h *Handler) HabitsGetAll(rw http.ResponseWriter, r *http.Request) {
	params := map[string]string{
		"startDate": r.URL.Query().Get("startDate"),
		"endDate":   r.URL.Query().Get("endDate"),
	}
	hs, err := h.App.HabitsGetAll(params)
	JsonCheckError(rw, http.StatusOK, envelope{"habits": types.FromHabitToTransportSlice(hs)}, err)
}
