package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type envelope map[string]interface{}

func (a *application) writeJson(rw http.ResponseWriter, status int, data envelope,
	headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')
	for k, v := range headers {
		rw.Header()[k] = v
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(js)

	return nil
}

func (a *application) validateDateParam(r *http.Request) (string, error) {
	// TODO validate that we got a usable date string
	date := chi.URLParam(r, "date")
	return date, nil
}
