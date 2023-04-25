package web

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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
	h.render(w, http.StatusOK, "view.html.tmpl", data)
}

func (h *Handler) habitCreateForm(w http.ResponseWriter, r *http.Request) {
	data := h.newTemplateData(r)
	data.Form = habitCreateForm{
		Expires: 1, // Default value.
	}
	h.render(w, http.StatusOK, "create.html.tmpl", data)
}

func (h *Handler) habitCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}
	// Date picker on the fron-ent that results in a yyyy-mm-dd
	// Radial buttons on the front-end that correspond to catefory IDs

	// The rest of the inputs for the habit log
	// - Success
	// - Note (description)
	// - Origin (automatically manual on the form)

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

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	data := h.newTemplateData(r)
	if !form.Valid() {
		data.Form = form
		// 422 - unprocessable entity represents a validation error.
		h.render(w, http.StatusUnprocessableEntity, "create.html.tmpl", data)
		return
	}

	d, err := h.App.DayGetByDate(time.Now())
	if err != nil {
		h.serverError(w, err)
		return
	}

	habit, err := h.App.HabitCreate(d, &entity.Habit{
		CategoryID: 1,
		Logs: []*entity.HabitLog{
			{
				Success: true,
				Origin:  entity.Manual,
				Note:    fmt.Sprintf("title: %s\n content: %s\n expires: %d\n", form.Title, form.Content, expires),
			},
		},
	})
	if err != nil {
		h.serverError(w, err)
		return
	}

	data.Habit = habit
	h.render(w, http.StatusCreated, "view.html.tmpl", data)
}
