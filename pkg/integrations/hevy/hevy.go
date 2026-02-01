package hevy

import "net/http"

type API struct {
	WorkoutEvents *WorkoutEvents
	Workouts      *Workouts
}

func New(apiKey string) *API {
	client := &Client{
		Client: http.DefaultClient,
		apiKey: apiKey,
	}
	return &API{
		WorkoutEvents: &WorkoutEvents{client: client},
		Workouts:      &Workouts{client: client},
	}
}
