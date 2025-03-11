package transport

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func errServer(err error) error {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	slog.Error(trace)
	return echo.NewHTTPError(http.StatusInternalServerError)
}

/*
func errClientError(err error) error {
	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
}
*/
