package main

import (
	"fmt"
	"net/http"
	"time"
)

func (a *application) miscHandler(rw http.ResponseWriter, r *http.Request) {
	t, err := a.collectors.Fit.DayLog("2021-05-25")
	if err != nil {
		a.badRequestResponse(rw, r, err)
	}

	msg := t
	e := envelope{
		"success": true,
		"message": msg,
	}
	a.writeJSON(rw, http.StatusOK, e, nil)
}

func (a *application) seedDataHandler(rw http.ResponseWriter, r *http.Request) {
	t := time.Now()

	err := a.collectors.Sleep.GetRecordsFromDayZero(t)
	if err != nil {
		a.errorResponse(rw, r, http.StatusNotFound, err)
		return
	}

	go a.collectors.Work.GetRecordsFromDayZero(t)
	if err != nil {
		a.errorResponse(rw, r, http.StatusNotFound, err)
		return
	}

	a.collectors.Sleep.BuildHabitsFromSleepRecords()
	a.collectors.Work.BuildHabitsFromWorkRecords()
	if err != nil {
		a.errorResponse(rw, r, http.StatusNotFound, err)
		return
	}

	ts := t.Format("2006-01-02")
	e := envelope{
		"message": fmt.Sprintf("Sleep records and habits build until %v", ts),
		"success": true,
	}
	a.writeJSON(rw, http.StatusOK, e, nil)
}
