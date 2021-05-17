package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/danielcosme/curious-ape/internal/data"
)

func (a *application) createFoodHabitHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Create food habit record")
}

func (a *application) showFoodHabitHandler(rw http.ResponseWriter, r *http.Request) {
	date, err := a.validateDateParam(r)
	if err != nil {
		http.NotFound(rw, r)
		return
	}

	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		a.logger.Println(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	habit := data.FoodHabit{
		ID:    1,
		State: true,
		Date:  dateTime,
		Tags:  []string{"lion", "16/8"},
	}

	err = a.writeJson(rw, http.StatusOK, habit, nil)
	if err != nil {
		a.logger.Println(err)
		http.Error(rw, "Internal Error", http.StatusInternalServerError)
	}
}
