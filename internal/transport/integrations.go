package transport

import (
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

func (t *Transport) Oauth2Success(c echo.Context) error {
	provider := c.Param("provider")
	code := c.QueryParam("code")
	err := t.App.Oauth2Success(provider, code)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/integrations")
}
