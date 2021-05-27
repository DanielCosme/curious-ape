package main

import (
	"net/http"
)

func (a *application) authorizeFitbitHandler(rw http.ResponseWriter, r *http.Request) {
	reqURL := a.collectors.Sleep.AuthorizationURI()
	http.Redirect(rw, r, reqURL, http.StatusFound)
}

func (a *application) successFitbitHandler(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		a.serverErrorResponse(rw, r, err)
		return
	}
	code := r.Form.Get("code") // Valid only for 10 minutes.

	payload, err := a.collectors.Sleep.Auth.ExchangeCodeForToken(code)
	if err != nil {
		a.serverErrorResponse(rw, r, err)
		return
	}

	err = a.models.Tokens.Update(payload)
	if err != nil {
		a.serverErrorResponse(rw, r, err)
	}

	a.writeJSON(rw, 200, envelope{"message": payload}, nil)
}
