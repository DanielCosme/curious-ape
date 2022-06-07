package router

import (
	"github.com/danielcosme/curious-ape/internal/api/types"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/rest"
	"net/http"
	"time"
)

func (h *Handler) HabitCategories(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		categories, err := h.App.HabitsGetCategories()
		JsonCheckError(rw, r, http.StatusOK, &envelope{"categories": categories}, err)
	}
}

func (h *Handler) Habits(rw http.ResponseWriter, r *http.Request) {
	var data *types.HabitTransport
	habit := r.Context().Value("habit").(*entity.Habit)

	switch r.Method {
	case http.MethodGet:
		if habit != nil {
			rest.JSONStatusOk(rw, envelope{"habit": types.FromHabitToTransport(habit)})
		} else {
			hs, err := h.HabitsGetAll()
			JsonCheckError(rw, r, http.StatusOK, envelope{"habits": hs}, err)
		}
	case http.MethodPost:
		day := r.Context().Value("day").(*entity.Day)

		err := rest.ReadJSON(r, &data)
		if err != nil {
			rest.ErrInternalServer(rw)
			return
		}

		newHabit, err := h.App.HabitCreate(day, data.ToHabit())
		JsonCheckError(rw, r, http.StatusCreated, envelope{"habit": types.FromHabitToTransport(newHabit)}, err)
	case http.MethodPut:
		err := rest.ReadJSON(r, &data)
		if err != nil {
			rest.ErrInternalServer(rw)
			return
		}

		habitUpdated, err := h.App.HabitFullUpdate(habit, data.ToHabit())
		JsonCheckError(rw, r, http.StatusOK, envelope{"habit": types.FromHabitToTransport(habitUpdated)}, err)
	case http.MethodDelete:
		err := h.App.HabitDelete(habit)
		JsonCheckError(rw, r, http.StatusOK, nil, err)
	default:
		rest.ErrNotAllowed(rw)
	}
}

func (h *Handler) HabitsGetAll() ([]*types.HabitTransport, error) {
	hs, err := h.App.HabitsGetAll(time.Now(), time.Now())
	if err != nil {
		return nil, err
	}

	habits := []*types.HabitTransport{}
	for _, habit := range hs {
		habits = append(habits, types.FromHabitToTransport(habit))
	}

	return habits, nil
}
