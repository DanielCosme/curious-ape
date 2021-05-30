package main

import (
	"fmt"
	"net/http"
	"time"
)

func (a *application) seedDataHandler(rw http.ResponseWriter, r *http.Request) {
	// build the app from scratch
	t := time.Now()
	err := a.collectors.FromDayZero(t)
	ts := t.Format("2006-01-02")
	if err != nil {
		a.serverErrorResponse(rw, r, err)
	}

	e := envelope{
		"message": fmt.Sprintf("Sleep records and habits build until %v", ts),
		"success": true,
	}
	a.writeJSON(rw, http.StatusOK, e, nil)
}
