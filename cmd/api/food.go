package main

import (
	"fmt"
	"net/http"
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

	fmt.Fprintf(rw, "The date is %s", date)
}
