package transport

import (
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (t *Transport) integrationsForm(c echo.Context) error {
	integrations, err := t.App.IntegrationsGet()
	if err != nil {
		return err
	}
	td := t.newTemplateData(c.Request())
	td.Integrations = integrations
	return c.Render(http.StatusOK, pageIntegrations, td)
}

func (t *Transport) oauth2Success(c echo.Context) error {
	provider := c.Param("provider")
	code := c.QueryParam("code")
	err := t.App.Oauth2Success(provider, code)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/integrations")
}

func (t *Transport) sync(c echo.Context) error {
	d, err := core.DateFromISO8601(c.Param("date"))
	if err != nil {
		return err
	}
	day, err := t.App.SyncDay(d)
	if err != nil {
		return err
	}

	td := t.newTemplateData(c.Request())
	td.Day = &formatDays([]core.Day{day})[0]
	return c.Render(http.StatusOK, partialDayRow, td.Day)
}
