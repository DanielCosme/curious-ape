package httprest

import (
	"net/http"
)

func (a *API) Habits(rw http.ResponseWriter, r *http.Request) {
	// Decode body if anything


	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
	case http.MethodPut:
	case http.MethodDelete:
	}
}
