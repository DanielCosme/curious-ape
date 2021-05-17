package main

import (
	"fmt"
	"net/http"
)

// Show api information
func (a *application) healthcheckerHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "status: available")
	fmt.Fprintf(rw, "environment: %s\n", a.config.env)
	fmt.Fprintf(rw, "version: %s\n", version)
}
