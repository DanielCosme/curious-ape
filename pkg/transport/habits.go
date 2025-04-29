package transport

import (
	"github.com/danielcosme/curious-ape/pkg/application"
	"net/http"
	"time"
)

type CreateHabitParams struct {
	DateParam string                 `json:"date"`
	State     application.HabitState `json:"state"`
	Type      application.HabitType  `json:"type"`
	Date      time.Time              `json:"-"`
}

func (hr *CreateHabitParams) Validate(v *Validator) {
	if d, err := time.Parse("2006-01-02", hr.DateParam); err != nil {
		v.Add("date", err.Error())
	} else {
		hr.Date = d
	}
	v.Check(application.StateValid(hr.State), "state", invalid(string(hr.State)))
	v.Check(application.HabitTypeValid(hr.Type), "type", invalid(string(hr.Type)))
	v.Message("invalid habit parameters")
}

func (t *Transport) HandlerHabitsUpsert(w http.ResponseWriter, r *http.Request) {
	var params CreateHabitParams
	if err := Bind(r.Body, &params); err != nil {
		JSONError(w, err, http.StatusBadRequest)
		return
	}
	habit, _ := t.app.HabitUpsert(params.Date, params.State, params.Type)
	JSON(w, http.StatusCreated, envelope{"habit": habit}, nil)
}

func (t *Transport) HandlerHabitTypes(w http.ResponseWriter, r *http.Request) {
	types := []application.HabitType{
		application.HabitTypeWake,
		application.HabitTypeWorkout,
		application.HabitTypeWork,
		application.HabitTypeEat,
	}
	JSONOK(w, envelope{"types": types}, nil)
}
