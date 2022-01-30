package stdmux

//
// import (
// 	"fmt"
// 	"net/http"
// )
//
// func (a *main.application) logError(r *http.Request, err error) {
// 	// a.logger.Println(err)
// 	fmt.Println(err.Error())
// }
//
// // json formated error messages.
// func (a *main.application) errorResponse(rw http.ResponseWriter, r *http.Request, status int, message interface{}) {
// 	e := envelope{
// 		"error":   message,
// 		"success": false,
// 	}
//
// 	err := a.writeJSON(rw, status, e, nil)
// 	if err != nil {
// 		a.logError(r, err)
// 		rw.WriteHeader(500)
// 	}
// }
//
// func (a *main.application) serverErrorResponse(rw http.ResponseWriter, r *http.Request,
// 	err error) {
// 	a.logError(r, err)
// 	msg := "unexpected problem encountered on the server, unable to process request"
// 	s := http.StatusInternalServerError
// 	a.errorResponse(rw, r, s, msg)
// }
//
// func (a *main.application) notFoundResponse(rw http.ResponseWriter, r *http.Request) {
// 	msg := "resource not found"
// 	a.errorResponse(rw, r, http.StatusNotFound, msg)
// }
//
// func (a *main.application) methodNotAllowedResponse(rw http.ResponseWriter, r *http.Request) {
// 	msg := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
// 	a.errorResponse(rw, r, http.StatusMethodNotAllowed, msg)
// }
//
// func (a *main.application) badRequestResponse(rw http.ResponseWriter, r *http.Request, err error) {
// 	a.errorResponse(rw, r, http.StatusBadRequest, err.Error())
// }
//
// func (a *main.application) failedValidationResponse(rw http.ResponseWriter, r *http.Request,
// 	errors map[string]string) {
// 	a.errorResponse(rw, r, http.StatusUnprocessableEntity, errors)
// }
//
// func (a *main.application) rateLimitExceededResponse(rw http.ResponseWriter, r *http.Request) {
// 	message := "rate limit exceeded"
// 	a.errorResponse(rw, r, http.StatusTooManyRequests, message)
// }
//
// func (a *main.application) unauthorizedResponse(rw http.ResponseWriter, r *http.Request) {
// 	rw.Header().Add("WWW-Authenticate", "basic")
// 	a.errorResponse(rw, r, http.StatusUnauthorized, "unauthorized")
// }
//
// func (a *main.application) invalidCredentialsResponse(rw http.ResponseWriter, r *http.Request) {
// 	rw.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
// 	a.errorResponse(rw, r, http.StatusUnauthorized, "invalid credentials")
// }
//
