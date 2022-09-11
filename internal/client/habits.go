package client

import (
	"github.com/danielcosme/curious-ape/internal/api/types"
	"net/http"
)

const DATE_FORMAT = "2002-01-02T00:00:00Z"

type HabitsService struct {
	C *Service
}

type habitsEnvelope struct {
	Habits []types.HabitTransport `json:"habits"`
}

func (h *HabitsService) List() ([]types.HabitTransport, error) {
	var habitsEnvelope habitsEnvelope
	if err := h.C.Call(http.MethodGet, "/habits", nil, &habitsEnvelope, nil); err != nil {
		return nil, err
	}
	return habitsEnvelope.Habits, nil
}
