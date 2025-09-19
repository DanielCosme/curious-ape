package api

import (
	"errors"
	"net/http"

	"github.com/danielcosme/curious-ape/pkg/persistence"
	"github.com/danielcosme/curious-ape/pkg/validator"
	"github.com/danielcosme/curious-ape/views"
	"github.com/labstack/echo/v4"
)

type userLoginForm struct {
	Username string
	Password string
	validator.Validator
}

func (api *API) getLogin(c echo.Context) error {
	if api.IsAuthenticated(c.Request()) {
		return c.Redirect(http.StatusFound, "/")
	}
	return renderOK(c, views.Login(views.State{Version: api.Version}))
}

func (api *API) loginPost(c echo.Context) error {
	form := userLoginForm{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
	}

	if !form.Valid() {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid form"))
	}

	id, err := api.App.Authenticate(form.Username, form.Password)
	if err != nil {
		if errors.Is(err, persistence.ErrInvalidCredentials) {
			form.AddNonFieldError("username or password is incorrect")
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	err = api.SessionManager.RenewToken(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	api.SessionManager.Put(c.Request().Context(), string(ctxKeyAuthenticatedUserID), id)
	return c.Redirect(http.StatusFound, "/")
}

func (api *API) logout(c echo.Context) error {
	if err := api.SessionManager.RenewToken(c.Request().Context()); err != nil {
		return errServer(err)
	}

	api.SessionManager.Remove(c.Request().Context(), string(ctxKeyAuthenticatedUserID))
	return c.NoContent(http.StatusOK)
}
