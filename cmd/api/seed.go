package main

import (
	"fmt"
	"net/http"
	"time"
)

func (a *application) dayGetHandler(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	d := r.Form.Get("day")

	msg := "no params"
	var suc bool
	prov := r.Form.Get("provider")
	switch prov {
	case "work":
		err := a.collectors.Work.GetRecord(d)
		if err != nil {
			a.serverErrorResponse(rw, r, err)
			return
		}
		msg = "All Good"
		suc = true
	case "fitness":
		err := a.collectors.Fit.GetRecord(d)
		if err != nil {
			a.serverErrorResponse(rw, r, err)
			return
		}
		msg = "All Good"
		suc = true
	case "habits":
		err := a.collectors.InitializeDayHabit()
		if err != nil {
			a.serverErrorResponse(rw, r, err)
			return
		}
		msg = "All Good"
		suc = true
	}

	e := envelope{"success": suc, "msg": msg}
	a.writeJSON(rw, http.StatusOK, e, nil)
}

func (a *application) miscHandler(rw http.ResponseWriter, r *http.Request) {

	err := a.collectors.Sleep.BuildHabitsFromSleepRecords()
	if err != nil {
		a.errorResponse(rw, r, 400, err)
		return
	}

	a.collectors.Work.BuildHabitsFromWorkRecords()
	if err != nil {
		a.errorResponse(rw, r, 400, err)
		return
	}
	err = a.collectors.Fit.BuildHabitsFromFitnessRecords()
	if err != nil {
		a.errorResponse(rw, r, 400, err)
		return
	}

	c := "all good"
	msg := c
	e := envelope{
		"success": true,
		"message": msg,
	}
	a.writeJSON(rw, http.StatusOK, e, nil)
}

func (a *application) seedDataHandler(rw http.ResponseWriter, r *http.Request) {
	t := time.Now()

	a.collectors.AllHabitsInit()
	err := a.collectors.Sleep.GetRecordsFromDayZero(t)
	if err != nil {
		a.errorResponse(rw, r, http.StatusNotFound, err)
		return
	}

	err = a.collectors.Fit.GetRecordsFromDayZero(t)
	if err != nil {
		a.errorResponse(rw, r, http.StatusNotFound, err.Error())
		return
	}

	go a.collectors.Work.GetRecordsFromDayZero(t)
	if err != nil {
		a.errorResponse(rw, r, http.StatusNotFound, err)
		return
	}

	ts := t.Format("2006-01-02")
	e := envelope{
		"message": fmt.Sprintf("Records and habits build until %v", ts),
		"success": true,
	}
	a.writeJSON(rw, http.StatusOK, e, nil)
}
