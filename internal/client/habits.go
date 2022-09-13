package client

import (
	"fmt"
	"github.com/danielcosme/curious-ape/internal/api/types"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"net/http"
	"net/url"
	"time"
)

const DATE_FORMAT = "2002-01-02T00:00:00Z"

type HabitsService struct {
	C *Service
}

type habitsEnvelope struct {
	Habits []types.HabitTransport `json:"habits"`
}

// TODO make this go away
type habitEnvelope struct {
	Habit *types.HabitTransport `json:"habit"`
}

type categoriesEnvelope struct {
	Categories []*entity.HabitCategory `json:"categories"`
}

func (h *HabitsService) List(startDate, endDate time.Time) ([]types.HabitTransport, error) {
	var habitsEnvelope habitsEnvelope
	p := url.Values{}
	p.Add("startDate", startDate.Format(entity.ISO8601))
	p.Add("endDate", endDate.Format(entity.ISO8601))

	if err := h.C.Call(http.MethodGet, "/habits", nil, &habitsEnvelope, p); err != nil {
		return nil, err
	}
	return habitsEnvelope.Habits, nil
}

func (h *HabitsService) Create(date time.Time, habit *types.HabitTransport) (*types.HabitTransport, error) {
	var habitEnvelope habitEnvelope
	path := fmt.Sprintf("/habits/date/%s", date.Format(entity.ISO8601))
	if err := h.C.Call(http.MethodPost, path, habit, &habitEnvelope, nil); err != nil {
		return nil, err
	}
	return habitEnvelope.Habit, nil
}

func (h *HabitsService) Categories() ([]*entity.HabitCategory, error) {
	var csEnvelope categoriesEnvelope
	if err := h.C.Call(http.MethodGet, "/habits/categories", nil, &csEnvelope, nil); err != nil {
		return nil, err
	}
	return csEnvelope.Categories, nil
}
