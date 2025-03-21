package transport

import (
	"github.com/danielcosme/curious-ape/pkg/database/gen/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type User struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

func (t *Transport) getUser(c echo.Context) error {
	u, ok := c.Request().Context().Value(ctxUser).(*models.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, User{
		Name: u.Username,
		Role: u.Role,
	})
}

type Info struct {
	Version string `json:"version"`
}

func (t *Transport) getVersion(c echo.Context) error {
	return c.JSON(http.StatusOK, Info{
		Version: t.Version,
	})
}
