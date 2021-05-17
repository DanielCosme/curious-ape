package main

import (
	"fmt"
	"net/http"

	"github.com/danielcosme/curious-ape/internal/data"
)

func (a *application) createFoodHabitHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Create food habit record")
}

func (a *application) showFoodHabitHandler(rw http.ResponseWriter, r *http.Request) {
	date, err := a.validateDateParam(r)
	if err != nil {
		a.errorResponse(rw, r, http.StatusBadRequest, err.Error())
		return
	}

	habit := data.FoodHabit{
		ID:    1,
		State: true,
		Date:  date,
		Tags:  []string{"lion", "16/8"},
	}

	err = a.writeJSON(rw, http.StatusOK, envelope{"foodHabit": habit}, nil)
	if err != nil {
		a.serverErrorResponse(rw, r, err)
	}
}
