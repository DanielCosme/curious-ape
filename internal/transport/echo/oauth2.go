package echo

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) Oauth2Success(c echo.Context) error {
	provider := c.Param("provider")
	code := c.QueryParam("code")

	err := h.App.Oauth2Success(provider, code)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (h *Handler) Oauth2Connect(c echo.Context) error {
	provider := c.Param("provider")

	url, err := h.App.Oauth2ConnectProvider(provider)
	if err != nil {
		return err
	}

	c.Response().Header().Set("location", url)
	return c.NoContent(http.StatusTemporaryRedirect)
}
