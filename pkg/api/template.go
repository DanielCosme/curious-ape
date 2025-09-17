package api

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func render(c echo.Context, status int, co templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)
	err := co.Render(ctx(c), buf)
	if err != nil {
		return err
	}
	return c.HTMLBlob(status, buf.Bytes())
}

func renderOK(c echo.Context, co templ.Component) error {
	return render(c, http.StatusOK, co)
}
