package main

import (
	"net/http"
)

// Show api information
func (a *application) healthcheckerHandler(rw http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status": "available",
		"systemInfo": map[string]string{
			"environment": a.config.env,
			"version":     version,
		},
	}

	err := a.writeJSON(rw, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(rw, r, err)
	}
}
