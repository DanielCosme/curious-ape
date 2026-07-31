package web

import (
	"io"
	"log/slog"
	"net/http"
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
)

type Node interface {
	Render(w io.Writer) error
}

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
func Render(w http.ResponseWriter, node Node) {
	w.Header().Set("Content-Type", "text/html")
	err := node.Render(w)
	if err != nil {
		ErrInternalServer(err, w)
	}
}

func GetDayParams(r *http.Request) core.Date {
	err := r.ParseForm()
	if err != nil {
		panic("get day params: error parsing form: " + err.Error())
	}

	if r.Form.Get("date") == "" {
		return core.NewDate(time.Now())
	}
	date, err := core.NewDateFromISO8601(r.Form.Get("date"))
	if err != nil {
		slog.Error("cannot parse date", "err", err)
		panic(err)
	}
	return date
}
