package main

import (
	"net/http"
)

// Show api information
func (a *application) healthcheckerHandler(rw http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": a.config.env,
		"version":     version,
	}

	err := a.writeJson(rw, http.StatusOK, data, nil)
	if err != nil {
		a.logger.Println(err)
		http.Error(rw, "Server error", http.StatusInternalServerError)
		return
	}
}
