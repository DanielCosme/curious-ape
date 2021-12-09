package api

//
// import (
// 	"errors"
// 	"fmt"
// 	"github.com/danielcosme/curious-ape/internal/core/entity"
// 	"github.com/danielcosme/curious-ape/internal/models/pg"
// 	"net/http"
// 	"strconv"
// 	"strings"
//
// 	"github.com/danielcosme/curious-ape/internal/validator"
// 	"github.com/go-chi/chi/v5"
// )
//
// func (a *application) createHabitHandler(rw http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		State  string `json:"state"`
// 		Date   string `json:"date"`
// 		Origin string `json:"origin"`
// 		Type   string `json:"type"`
// 	}
//
// 	err := a.readJSON(rw, r, &input)
// 	if err != nil {
// 		a.badRequestResponse(rw, r, err)
// 		return
// 	}
//
// 	habit := &entity.Habit{
// 		Date:   input.Date,
// 		State:  input.State,
// 		Origin: input.Origin,
// 		Type:   input.Type,
// 	}
// 	v := validator.New()
// 	if pg.ValidateHabit(v, habit); !v.Valid() {
// 		a.failedValidationResponse(rw, r, v.Errors)
// 		return
// 	}
//
// 	err = a.models.Habits.UpdateOrCreate(habit)
// 	if err != nil {
// 		a.serverErrorResponse(rw, r, err)
// 		return
// 	}
//
// 	headers := make(http.Header)
// 	headers.Set("Location", fmt.Sprintf("/v1/habits/%d", habit.ID))
//
// 	err = a.writeJSON(rw, http.StatusCreated, envelope{"Habit": habit}, headers)
// 	if err != nil {
// 		a.serverErrorResponse(rw, r, err)
// 	}
// }
//
// func (a *application) listHabitsHandler(rw http.ResponseWriter, r *http.Request) {
// 	habits, err := a.models.Habits.GetAll()
// 	if err != nil {
// 		a.serverErrorResponse(rw, r, err)
// 		return
// 	}
//
// 	err = a.writeJSON(rw, http.StatusOK, envelope{"habits": habits}, nil)
// 	if err != nil {
// 		a.serverErrorResponse(rw, r, err)
// 	}
// }
//
// func (a *application) showHabitHandler(rw http.ResponseWriter, r *http.Request) {
// 	id := chi.URLParam(r, "id")
// 	idInt, err := strconv.Atoi(id)
// 	if err != nil {
// 		a.badRequestResponse(rw, r, err)
// 		return
// 	}
//
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
//
// 	err = a.writeJSON(rw, http.StatusOK, envelope{"Habit": habit}, nil)
// 	if err != nil {
// 		a.serverErrorResponse(rw, r, err)
// 	}
// }
//
// func (a *application) updateHabitHandler(rw http.ResponseWriter, r *http.Request) {
// 	id := chi.URLParam(r, "id")
// 	idInt, err := strconv.Atoi(id)
// 	if err != nil {
// 		a.badRequestResponse(rw, r, err)
// 		return
// 	}
//
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
//
// 	var input struct {
// 		State  string `json:"state"`
// 		Date   string `json:"date"`
// 		Origin string `json:"origin"`
// 		Type   string `json:"type"`
// 	}
// 	err = a.readJSON(rw, r, &input)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "body must not be empty") {
// 			a.errorResponse(rw, r, http.StatusBadRequest, err.Error())
// 			return
// 		}
// 		a.serverErrorResponse(rw, r, err)
// 		return
// 	}
//
// 	if habit.Type != input.Type {
// 		a.badRequestResponse(rw, r, fmt.Errorf("habits must be the same type"))
// 		return
// 	}
//
// 	habit.State = input.State
// 	habit.Date = input.Date
// 	habit.Type = input.Type
// 	habit.Origin = input.Origin
// 	v := validator.New()
// 	if pg.ValidateHabit(v, habit); !v.Valid() {
// 		a.failedValidationResponse(rw, r, v.Errors)
// 		return
// 	}
//
// 	err = a.models.Habits.Update(habit)
// 	if err != nil {
// 		a.serverErrorResponse(rw, r, err)
// 		return
// 	}
//
// 	err = a.writeJSON(rw, http.StatusOK, envelope{"Habit": habit}, nil)
// 	if err != nil {
// 		a.serverErrorResponse(rw, r, err)
// 		return
// 	}
// }
//
// func (a *application) deleteHabitHandler(rw http.ResponseWriter, r *http.Request) {
// 	id, err := strconv.Atoi(chi.URLParam(r, "id"))
// 	if err != nil {
// 		a.badRequestResponse(rw, r, err)
// 		return
// 	}
//
// 	err = a.models.Habits.Delete(id)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, errors.ErrRecordNotFound):
// 			a.notFoundResponse(rw, r)
// 		default:
// 			a.serverErrorResponse(rw, r, err)
// 		}
// 		return
// 	}
//
// 	err = a.writeJSON(rw, http.StatusOK, envelope{"message": "food habit log successfully deleted"}, nil)
// 	if err != nil {
// 		a.serverErrorResponse(rw, r, err)
// 	}
// }
//
