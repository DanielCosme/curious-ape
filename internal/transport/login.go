package transport

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/danielcosme/curious-ape/internal/database"

	"github.com/danielcosme/curious-ape/internal/validator"
)

type userLoginForm struct {
	Email    string
	Password string
	validator.Validator
}

func (t *Transport) loginForm(c echo.Context) error {
	data := t.newTemplateData(c.Request())
	data.Form = userLoginForm{}
	return c.Render(http.StatusOK, pageLogin, data)
}

func (t *Transport) loginPost(c echo.Context) error {
	form := userLoginForm{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	if !form.Valid() {
		data := t.newTemplateData(c.Request())
		data.Form = form
		return c.Render(http.StatusUnprocessableEntity, pageLogin, data)
	}

	id, err := t.App.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, database.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := t.newTemplateData(c.Request())
			data.Form = form
			return c.Render(http.StatusUnprocessableEntity, pageLogin, data)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	err = t.SessionManager.RenewToken(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	t.SessionManager.Put(c.Request().Context(), string(ctxKeyAuthenticatedUserID), id)
	return c.Redirect(http.StatusSeeOther, "/")
}

func (t *Transport) logout(c echo.Context) error {
	if err := t.SessionManager.RenewToken(c.Request().Context()); err != nil {
		return errServer(err)
	}

	t.SessionManager.Remove(c.Request().Context(), "authenticatedUserID")
	t.SessionManager.Put(c.Request().Context(), "flash", "You've been logged out successfully!")
	return c.Redirect(http.StatusSeeOther, "/")
}
