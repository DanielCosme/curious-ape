package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *application) validateDateParam(r *http.Request) (string, error) {
	// TODO validate that we got a usable date string
	date := chi.URLParam(r, "date")
	return date, nil
}
