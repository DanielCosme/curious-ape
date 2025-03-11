package transport

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

func (t *Transport) home(c echo.Context) error {
	w := c.Response().Writer
	_, err := fmt.Fprintf(w, "<h1>Hello World!</h1>")
	return errServer(err)
}
