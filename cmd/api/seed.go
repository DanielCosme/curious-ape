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
		a.errorResponse(rw, r, http.StatusNotFound, err)
	}

	e := envelope{
		"message": fmt.Sprintf("Sleep records and habits build until %v", ts),
		"success": true,
	}
	a.writeJSON(rw, http.StatusOK, e, nil)
}

func (a *application) miscHandler(rw http.ResponseWriter, r *http.Request) {
	err := a.collectors.GetLog("2021-05-26")
	if err != nil {
		a.errorResponse(rw, r, http.StatusNotFound, err.Error())
		return
	}

	msg := "yes"
	e := envelope{
		"message": msg,
		"success": true,
	}
	a.writeJSON(rw, http.StatusOK, e, nil)
}
