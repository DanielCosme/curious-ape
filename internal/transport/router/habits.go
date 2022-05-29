package router

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/transport/types"
	"github.com/danielcosme/curious-ape/rest"
	"net/http"
	"time"
)

func (h *Handler) HabitCategories(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		categories, err := h.App.HabitsGetCategories()
		JsonCheckError(rw, r, http.StatusOK, &rest.E{"categories": categories}, err)
	}
}

func (h *Handler) Habits(rw http.ResponseWriter, r *http.Request) {
	var data *types.HabitTransport
	habit := r.Context().Value("habit").(*entity.Habit)

	switch r.Method {
	case "GET":
		if habit != nil {
			rest.JSONStatusOk(rw, &rest.E{"habit": types.FromHabitToTransport(habit)})
		} else {
			hs, err := h.HabitsGetAll()
			JsonCheckError(rw, r, http.StatusOK, &rest.E{"habits": hs}, err)
		}
	case "POST":
		day := r.Context().Value("day").(*entity.Day)

		err := rest.ReadJSON(r, &data)
		if err != nil {
			rest.ErrInternalServer(rw, r)
			return
		}

		newHabit, err := h.App.HabitCreate(day, data.ToHabit())
		JsonCheckError(rw, r, http.StatusCreated, &rest.E{"habit": newHabit}, err)
	case "PUT":
		err := rest.ReadJSON(r, &data)
		if err != nil {
			rest.ErrInternalServer(rw, r)
			return
		}

		habitUpdated, err := h.App.HabitFullUpdate(habit, data.ToHabit())
		JsonCheckError(rw, r, http.StatusOK, &rest.E{"habit": habitUpdated}, err)
	case "DELETE":
		err := h.App.HabitDelete(habit)
		JsonCheckError(rw, r, http.StatusOK, nil, err)
	default:
		rest.ErrMethodNotSupported(rw, r)
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

// 	habit, err := a.models.Habits.Get(idInt)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, errors.ErrRecordNotFound):
// 			a.notFoundResponse(rw, r)
// 		default:
// 			a.serverErrorResponse(rw, r, err)
// 		}
// 		return
// 	}