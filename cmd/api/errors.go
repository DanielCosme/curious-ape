package main

import (
	"fmt"
	"net/http"
)

func (a *application) logError(r *http.Request, err error) {
	a.logger.Println(err)
}

// json formated error messages.
func (a *application) errorResponse(rw http.ResponseWriter, r *http.Request, status int,
	message interface{}) {
	e := envelope{"error": message}

	err := a.writeJSON(rw, status, e, nil)
	if err != nil {
		a.logError(r, err)
		rw.WriteHeader(500)
	}
}

func (a *application) serverErrorResponse(rw http.ResponseWriter, r *http.Request,
	err error) {
	a.logError(r, err)
	msg := "unexpected problem encountered on the server, unable to process request"
	s := http.StatusInternalServerError
	a.errorResponse(rw, r, s, msg)
}

func (a *application) notFoundResponse(rw http.ResponseWriter, r *http.Request) {
	msg := "resource not found"
	a.errorResponse(rw, r, http.StatusNotFound, msg)
}

func (a *application) methodNotAllowedResponse(rw http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	a.errorResponse(rw, r, http.StatusMethodNotAllowed, msg)
}
