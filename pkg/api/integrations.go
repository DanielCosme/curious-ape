package api

import (
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (api *API) integrationsGetAll(c echo.Context) error {
	integrations, err := api.App.IntegrationsGetAll()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, integrations)
}

func (api *API) integrationsGet(c echo.Context) error {
	provider := c.Param("provider")
	integration, err := api.App.IntegrationsGet(core.Integration(provider))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, integration)
}

func (api *API) oauth2Success(c echo.Context) error {
	provider := c.Param("provider")
	code := c.QueryParam("code")
	err := api.App.Oauth2Success(provider, code)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/integrations")
}
