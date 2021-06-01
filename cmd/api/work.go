package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *application) showWorkRecordHandler(rw http.ResponseWriter, r *http.Request) {
	date := chi.URLParam(r, "date")
	// TODO validate date input, for now it sends a 404 if not valid.

	record, err := a.models.WorkRecords.Get(date)
	if err != nil {
		a.notFoundResponse(rw, r)
		return
	}

	err = a.writeJSON(rw, http.StatusOK, envelope{"sleepRecord": record}, nil)
	if err != nil {
		a.serverErrorResponse(rw, r, err)
	}
}

func (a *application) listWorkRecordsHandler(rw http.ResponseWriter, r *http.Request) {
	record, err := a.models.WorkRecords.GetAll()
	if err != nil {
		a.notFoundResponse(rw, r)
		return
	}

	err = a.writeJSON(rw, http.StatusOK, envelope{"sleepRecords": record}, nil)
	if err != nil {
		a.serverErrorResponse(rw, r, err)
	}
}
