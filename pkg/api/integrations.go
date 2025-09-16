package api

import (
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (t *Transport) integrationsGetAll(c echo.Context) error {
	integrations, err := t.App.IntegrationsGetAll()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, integrations)
}

func (t *Transport) integrationsGet(c echo.Context) error {
	provider := c.Param("provider")
	integration, err := t.App.IntegrationsGet(core.Integration(provider))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, integration)
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
