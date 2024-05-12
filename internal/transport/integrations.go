package transport

import (
	"fmt"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (t *Transport) integrationsForm(c echo.Context) error {
	data := t.newTemplateData(c.Request())
	data.Fitbit = Integration{State: "Not Authenticated"}
	return c.Render(http.StatusOK, pageIntegrations, data)
}

func (t *Transport) oauth2Connect(c echo.Context) error {
	uri, err := t.App.Oauth2ConnectProvider(core.IntegrationFitbit)
	if err != nil {
		return err
	}

	t.App.Log.Info(uri)
	s := fmt.Sprintf(`<a href="%s"><button>Redirect</button></a>`, uri)
	return c.HTML(http.StatusOK, s)
}

func (t *Transport) Oauth2Success(c echo.Context) error {
	provider := c.Param("provider")
	code := c.QueryParam("code")

	err := t.App.Oauth2Success(provider, code)
	if err != nil {
		return err
	}
	t.App.Log.Info("Authentication successful", "provider", provider, "code", code)
	return c.Redirect(http.StatusSeeOther, "/integrations")
}
