package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/validator"
)

type newHabitForm struct {
	Date         time.Time
	CategoryCode string
	Success      bool
	Origin       entity.DataSource
	Note         string
	IsAutomated  bool
	validator.Validator
}

func (h *Handler) habit(w http.ResponseWriter, r *http.Request) {
	habit := r.Context().Value("habit").(*entity.Habit)

	data := h.newTemplateData(r)
	data.Habit = habit
	h.render(w, http.StatusOK, "view.gohtml", data)
}

func (h *Handler) newHabitForm(w http.ResponseWriter, r *http.Request) {
	data := h.newTemplateData(r)
	data.Form = newHabitForm{}
	h.render(w, http.StatusOK, "new_habit.gohtml", data)
}

func (h *Handler) newHabitPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	// Date: form
	// Success: form
	// Category code: form
	// - "food"
	// - "wake_up"
	// - "fitness"
	// - "deep_work"
	// - "custom"

	// HARDCODED
	// Origin: web
	// IsAutomated: False
	// Note: Empty

	form := newHabitForm{
		Origin:      entity.Manual,
		IsAutomated: false,
		// Title:   r.PostForm.Get("title"),
		// Content: r.PostForm.Get("content"),
		// Expires: expires,
	}

	// form.Check(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	// form.Check(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	// form.Check(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	// form.Check(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	data := h.newTemplateData(r)
	if !form.Valid() {
		data.Form = form
		// 422 - unprocessable entity represents a validation error.
		h.render(w, http.StatusUnprocessableEntity, "new_habit.gohtml", data)
		return
	}

	params := &application.NewHabitParams{
		Date:         time.Now(),
		CategoryCode: entity.HabitTypeFood.Str(),
		Success:      true,
		Origin:       entity.Manual,
		Note:         "",
		IsAutomated:  false,
	}
	habit, err := h.App.HabitUpsert(params)
	if err != nil {
		h.serverError(w, err)
		return
	}
	h.App.Session.Put(r.Context(), "flash", "Habit successfully created!")

	data.Habit = habit
	data.Flash = h.App.Session.PopString(r.Context(), "flash")
	h.render(w, http.StatusCreated, "view.gohtml", data)
}

func (h *Handler) newHabitLogPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	success, err := strconv.ParseBool(r.Form.Get("success"))
	if err != nil {
		h.serverError(w, err)
		return
	}
	dt, err := time.Parse(time.DateOnly, r.Form.Get("date"))
	if err != nil {
		h.serverError(w, err)
		return
	}

	params := &application.NewHabitParams{
		Date:         dt,
		CategoryCode: r.Form.Get("category"),
		Success:      success,
		Origin:       entity.Manual,
		IsAutomated:  false,
	}
	habit, err := h.App.HabitUpsert(params)
	if err != nil {
		h.serverError(w, err)
		return
	}
	day, err := h.App.DayGetByID(habit.DayID)
	if err != nil {
		h.serverError(w, err)
		return
	}

	td := h.newTemplateData(r)
	td.Day = &formatDays([]*entity.Day{day})[0]
	h.renderPartial(w, http.StatusOK, "day_row.gohtml", td)
}
