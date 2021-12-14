package rest

import "net/http"

type HTTPErr struct {
	Code     int    `json:"-"`
	Message  string `json:"message"`
	Internal error  `json:"-"`
}


func NewErr(rw http.ResponseWriter, code int, msg string) {
	rw.WriteHeader(code)
}
