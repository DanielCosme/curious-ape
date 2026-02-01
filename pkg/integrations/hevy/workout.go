package hevy

import (
	"net/http"
)

type Workouts struct {
	client *Client
}

type WorkoutCountEnvelope struct {
	WorkoutCount int `json:"workout_count"`
}

func (ws *Workouts) Count() (int, error) {
	var envelope WorkoutCountEnvelope
	err := ws.client.Call(http.MethodGet, "/v1/workouts/count", nil, &envelope)
	return envelope.WorkoutCount, err
}
