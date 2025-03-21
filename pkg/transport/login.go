package transport

import (
	"errors"
	"github.com/danielcosme/curious-ape/pkg/database"
	"github.com/danielcosme/curious-ape/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

type userLoginForm struct {
	Username string
	Password string
	validator.Validator
}

func (t *Transport) loginPost(c echo.Context) error {
	form := userLoginForm{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
	}

	if !form.Valid() {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid form"))
	}

	id, err := t.App.Authenticate(form.Username, form.Password)
	if err != nil {
		if errors.Is(err, database.ErrInvalidCredentials) {
			form.AddNonFieldError("username or password is incorrect")
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	err = t.SessionManager.RenewToken(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	t.SessionManager.Put(c.Request().Context(), string(ctxKeyAuthenticatedUserID), id)
	return c.NoContent(http.StatusOK)
}

func (t *Transport) logout(c echo.Context) error {
	if err := t.SessionManager.RenewToken(c.Request().Context()); err != nil {
		return errServer(err)
	}

	t.SessionManager.Remove(c.Request().Context(), string(ctxKeyAuthenticatedUserID))
	return c.NoContent(http.StatusOK)
}
