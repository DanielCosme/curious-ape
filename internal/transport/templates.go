package transport

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/danielcosme/curious-ape/internal/view"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func (t *Transport) RenderTempl(statusCode int, c echo.Context, tmpl templ.Component) error {
	c.Response().Writer.WriteHeader(statusCode)
	return tmpl.Render(c.Request().Context(), c.Response().Writer)
}

func (t *Transport) newTemplateData(r *http.Request) view.GlobalState {
	return view.GlobalState{
		Year:          fmt.Sprintf("%d", time.Now().Year()),
		Authenticated: t.IsAuthenticated(r),
		Version:       t.Version,
	}
}
