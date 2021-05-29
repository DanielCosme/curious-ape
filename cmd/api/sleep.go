package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *application) showSleepRecordHandler(rw http.ResponseWriter, r *http.Request) {
	date := chi.URLParam(r, "date")
	// TODO validate date input, for now it sends a 404 if not valid.

	record, err := a.models.SleepRecords.Get(date)
	if err != nil {
		a.notFoundResponse(rw, r)
		return
	}

	err = a.writeJSON(rw, http.StatusOK, envelope{"sleepRecord": record}, nil)
	if err != nil {
		a.serverErrorResponse(rw, r, err)
	}
}

func (a *application) listSleepRecordsHandler(rw http.ResponseWriter, r *http.Request) {
	record, err := a.models.SleepRecords.GetAll()
	if err != nil {
		a.notFoundResponse(rw, r)
		return
	}

	err = a.writeJSON(rw, http.StatusOK, envelope{"sleepRecords": record}, nil)
	if err != nil {
		a.serverErrorResponse(rw, r, err)
	}
}
