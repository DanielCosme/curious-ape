package httprest

import (
	"encoding/json"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"net/http"
)

func (a *API) HabitsGetAll(rw http.ResponseWriter, r *http.Request) {
	hs, err := a.App.Habits.Find(&entity.HabitQuery{})
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(hs)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Write(res)
	rw.WriteHeader(http.StatusOK)
}

func (a *API) HabitCreate(rw http.ResponseWriter, r *http.Request) {

}

func (a *API) HabitGet(rw http.ResponseWriter, r *http.Request) {

}

func (a *API) HabitUpdate(rw http.ResponseWriter, r *http.Request) {

}
