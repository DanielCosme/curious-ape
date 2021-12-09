package api

//
// import (
// 	"errors"
// 	"net/http"
// )
//
// func (a *main.application) authorizeFitbitHandler(rw http.ResponseWriter, r *http.Request) {
// 	reqURL := a.collectors.Sleep.AuthorizationURI()
// 	http.Redirect(rw, r, reqURL, http.StatusFound)
// }
//
// func (a *main.application) successFitbitHandler(rw http.ResponseWriter, r *http.Request) {
// 	err := r.ParseForm()
// 	if err != nil {
// 		a.serverErrorResponse(rw, r, err)
// 		return
// 	}
// 	code := r.Form.Get("code") // Valid only for 10 minutes.
// 	if code == "" {
// 		a.badRequestResponse(rw, r, errors.New("no authorization code provided"))
// 		return
// 	}
//
// 	payload, err := a.collectors.Sleep.Auth.ExchangeCodeForToken(code)
// 	if err != nil {
// 		a.serverErrorResponse(rw, r, err)
// 		return
// 	}
//
// 	err = a.models.Tokens.Update(payload)
// 	if err != nil {
// 		a.serverErrorResponse(rw, r, err)
// 	}
//
// 	e := envelope{
// 		"message": "authorization successful",
// 		"success": true,
// 	}
// 	a.writeJSON(rw, 200, e, nil)
// }
//
