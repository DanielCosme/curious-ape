package web

import (
	"io"
	"log/slog"
	"net/http"
)

type Node interface {
	Render(w io.Writer) error
}

func Redirect(w http.ResponseWriter, loc string) {
	w.Header().Set("Location", loc)
	w.WriteHeader(http.StatusSeeOther)
}

func Render(w http.ResponseWriter, node Node) {
	w.Header().Set("Content-Type", "text/html")
	err := node.Render(w)
	if err != nil {
		ErrInternalServer(err, w)
	}
}

func Err(code int, err error, w http.ResponseWriter) {
	slog.Error("Error", "err", err)
	http.Error(w, http.StatusText(code), code)
}

func ErrBadRequest(err error, w http.ResponseWriter) {
	slog.Error("Bad request", "err", err)
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func ErrInternalServer(err error, w http.ResponseWriter) {
	slog.Error("Internal server error", "err", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
