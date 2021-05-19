package main

import (
	"fmt"
	"net/http"

	"github.com/danielcosme/curious-ape/internal/data"
	"github.com/danielcosme/curious-ape/internal/validator"
	"github.com/go-chi/chi/v5"
)

func (a *application) createFoodHabitHandler(rw http.ResponseWriter, r *http.Request) {
	var input struct {
		State bool     `json:"state"`
		Date  string   `json:"date"`
		Tags  []string `json:"tags"`
	}

	err := a.readJSON(rw, r, &input)
	if err != nil {
		a.badRequestResponse(rw, r, err)
		return
	}

	habit := &data.FoodHabit{
		State: input.State,
		Date:  input.Date,
		Tags:  input.Tags,
	}

	v := validator.New()
	if data.ValidateFoodHabit(v, habit); !v.Valid() {
		a.failedValidationResponse(rw, r, v.Errors)
		return
	}

	err = a.models.FoodHabits.Insert(habit)
	if err != nil {
		a.serverErrorResponse(rw, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", habit.ID))

	err = a.writeJSON(rw, http.StatusCreated, envelope{"foodHabit": habit}, headers)
	if err != nil {
		a.serverErrorResponse(rw, r, err)
	}
}

func (a *application) showFoodHabitHandler(rw http.ResponseWriter, r *http.Request) {
	date := chi.URLParam(r, "date")
	habit := data.FoodHabit{
		ID:    1,
		State: true,
		Date:  date,
		Tags:  []string{"lion", "16/8"},
	}

	err := a.writeJSON(rw, http.StatusOK, envelope{"foodHabit": habit}, nil)
	if err != nil {
		a.serverErrorResponse(rw, r, err)
	}
}
