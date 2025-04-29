package transport

import (
	"fmt"
	"net/http"
)

func (t *Transport) HandlerDays(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		fmt.Fprintf(w, "GET, POST")
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (t *Transport) HandlerDaysMonth(w http.ResponseWriter, r *http.Request) {
	ds, _ := t.app.Days()
	JSONOK(w, envelope{"days": ds}, nil)
}
