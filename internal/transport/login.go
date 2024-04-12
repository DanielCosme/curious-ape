package transport

import (
	"errors"
	"net/http"

	"github.com/danielcosme/curious-ape/internal/core/database"
	"github.com/danielcosme/curious-ape/internal/validator"
)

type userLoginForm struct {
	Email    string
	Password string
	validator.Validator
}

func (h *Handler) loginForm(w http.ResponseWriter, r *http.Request) {
	data := h.newTemplateData(r)
	data.Form = userLoginForm{}
	h.render(w, http.StatusOK, "login.gohtml", data)
}

func (h *Handler) loginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.App.Log.Error(err)
		h.clientError(w, http.StatusBadRequest)
		return
	}

	form := userLoginForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	if !form.Valid() {
		data := h.newTemplateData(r)
		data.Form = form
		h.render(w, http.StatusUnprocessableEntity, "login.gohtml", data)
		return
	}

	id, err := h.App.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, database.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := h.newTemplateData(r)
			data.Form = form
			h.render(w, http.StatusUnprocessableEntity, "login.gohtml", data)
		} else {
			h.serverError(w, err)
		}
		return
	}

	err = h.SessionManager.RenewToken(r.Context())
	if err != nil {
		h.serverError(w, err)
		return
	}

	h.SessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if err := h.SessionManager.RenewToken(r.Context()); err != nil {
		h.serverError(w, err)
		return
	}

	h.SessionManager.Remove(r.Context(), "authenticatedUserID")
	h.SessionManager.Put(r.Context(), "flash", "You've been logged out sucessfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
