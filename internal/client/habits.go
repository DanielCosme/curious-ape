package client

import (
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
