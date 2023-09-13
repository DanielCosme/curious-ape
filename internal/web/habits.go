package web

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/validator"
)

type habitCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

func (h *Handler) habitView(w http.ResponseWriter, r *http.Request) {
	habit := r.Context().Value("habit").(*entity.Habit)

	data := h.newTemplateData(r)
	data.Habit = habit
	h.render(w, http.StatusOK, "view.gohtml", data)
}

func (h *Handler) habitCreateForm(w http.ResponseWriter, r *http.Request) {
	data := h.newTemplateData(r)
	data.Form = habitCreateForm{
		Expires: 1, // Default value.
	}
	h.render(w, http.StatusOK, "create.gohtml", data)
}

func (h *Handler) habitCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	form := habitCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	form.Check(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.Check(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.Check(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.Check(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	data := h.newTemplateData(r)
	if !form.Valid() {
		data.Form = form
		// 422 - unprocessable entity represents a validation error.
		h.render(w, http.StatusUnprocessableEntity, "create.gohtml", data)
		return
	}

	parms := &application.NewHabitParams{
		Date:         time.Now(),
		CategoryCode: entity.HabitTypeFood.Str(),
		Success:      true,
		Origin:       entity.Manual,
		Note:         fmt.Sprintf("title: %s\n content: %s\n expires: %d\n", form.Title, form.Content, expires),
		IsAutomated:  false,
	}
	habit, err := h.App.HabitUpsert(parms)
	if err != nil {
		h.serverError(w, err)
		return
	}
	h.App.Session.Put(r.Context(), "flash", "Habit successfully created!")

	data.Habit = habit
	data.Flash = h.App.Session.PopString(r.Context(), "flash")
	h.render(w, http.StatusCreated, "view.gohtml", data)
}
