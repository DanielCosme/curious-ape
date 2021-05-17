package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type envelope map[string]interface{}

func (a *application) writeJSON(rw http.ResponseWriter, status int, data envelope,
	headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for k, v := range headers {
		rw.Header()[k] = v
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(js)

	return nil
}

func (a *application) validateDateParam(r *http.Request) (time.Time, error) {
	date := chi.URLParam(r, "date")
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return dateTime, err
	}
	return dateTime, nil
}
