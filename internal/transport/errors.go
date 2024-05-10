package transport

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func (t *Transport) serverError(w http.ResponseWriter, err error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (t *Transport) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (t *Transport) notFound(w http.ResponseWriter) {
	t.clientError(w, http.StatusNotFound)
}

func errServer(err error) error {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	slog.Error(trace)
	return echo.NewHTTPError(http.StatusInternalServerError)
}

func errClientError() error {
	return echo.NewHTTPError(http.StatusBadRequest)
}

func errNotFound() error {
	return echo.NewHTTPError(http.StatusNotFound)
}
