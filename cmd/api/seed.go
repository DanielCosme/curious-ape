package main

import (
	"fmt"
	"net/http"
	"time"
)

func (a *application) dayGetHandler(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	d := r.Form.Get("day")

	prov := r.Form.Get("provider")
	var err error
	switch prov {
	case "work":
		err = a.collectors.Work.GetRecord(d)
	case "fitness":
		err = a.collectors.Fit.GetRecord(d)
	case "habits":
		err = a.collectors.InitializeDayHabit()
	}
	if err != nil {
		a.serverErrorResponse(rw, r, err)
		return
	}
	msg := "all good"
	suc := true

	a.writeJSON(rw, http.StatusOK, envelope{"success": suc, "msg": msg}, nil)
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

	// a.collectors.AllHabitsInit()
	// err := a.collectors.Sleep.GetRecordsFromDayZero(t)
	// if err != nil {
	// 	a.errorResponse(rw, r, http.StatusNotFound, err)
	// 	return
	// }

	// err = a.collectors.Fit.GetRecordsFromDayZero(t)
	// if err != nil {
	// 	a.errorResponse(rw, r, http.StatusNotFound, err.Error())
	// 	return
	// }

	err := a.collectors.Work.GetRecordsFromDayZero(t)
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
