package transport

import (
	"net/http"
	"strconv"
	"time"

	"github.com/danielcosme/curious-ape/internal/application"
	entity2 "github.com/danielcosme/curious-ape/internal/entity"

	"github.com/danielcosme/curious-ape/internal/validator"
)

type newHabitForm struct {
	Date         time.Time
	CategoryCode string
	Success      bool
	Origin       entity2.DataSource
	Note         string
	IsAutomated  bool
	validator.Validator
}

func (h *Transport) habit(w http.ResponseWriter, r *http.Request) {
	habit := r.Context().Value("habit").(*entity2.Habit)

	data := h.newTemplateData(r)
	data.Habit = habit
	h.render(w, http.StatusOK, "view.gohtml", data)
}

func (h *Transport) newHabitForm(w http.ResponseWriter, r *http.Request) {
	data := h.newTemplateData(r)
	data.Form = newHabitForm{}
	h.render(w, http.StatusOK, "new_habit.gohtml", data)
}

func (h *Transport) newHabitPost(w http.ResponseWriter, r *http.Request) {
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
	// Origin: transport
	// IsAutomated: False
	// Note: Empty

	form := newHabitForm{
		Origin:      entity2.Manual,
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
		CategoryCode: entity2.HabitTypeFood.Str(),
		Success:      true,
		Origin:       entity2.Manual,
		Note:         "",
		IsAutomated:  false,
	}
	habit, err := h.App.HabitUpsert(params)
	if err != nil {
		h.serverError(w, err)
		return
	}
	h.SessionManager.Put(r.Context(), "flash", "Habit successfully created!")

	data.Habit = habit
	data.Flash = h.SessionManager.PopString(r.Context(), "flash")
	h.render(w, http.StatusCreated, "view.gohtml", data)
}

func (h *Transport) newHabitLogPost(w http.ResponseWriter, r *http.Request) {
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
		Origin:       entity2.Manual,
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
	td.Day = &formatDays([]*entity2.Day{day})[0]
	h.renderPartial(w, http.StatusOK, "day_row.gohtml", td)
}
