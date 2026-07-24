package web

import (
	"log/slog"
	"net/http"
)

func Redirect(w http.ResponseWriter, loc string) {
	w.Header().Set("Location", loc)
	w.WriteHeader(http.StatusSeeOther)
}

func ErrInternalServer(err error, w http.ResponseWriter) {
	slog.Error("Internal server error", "err", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func Err(code int, err error, w http.ResponseWriter) {
	slog.Error("Error", "err", err)
	http.Error(w, http.StatusText(code), code)
}
