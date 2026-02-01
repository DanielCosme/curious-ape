package hevy

import (
	"net/http"
	"net/url"
	"time"
)

type WorkoutEvents struct {
	client *Client
}

type WorkoutEventsEnvelope struct {
	Page      int     `json:"page"`
	PageCount int     `json:"page_count"`
	Events    []Event `json:"events"`
}
type Set struct {
	Index           int     `json:"index"`
	Type            string  `json:"type"`
	WeightKg        float64 `json:"weight_kg"`
	Reps            int     `json:"reps"`
	DistanceMeters  any     `json:"distance_meters"`
	DurationSeconds any     `json:"duration_seconds"`
	Rpe             any     `json:"rpe"`
	CustomMetric    any     `json:"custom_metric"`
}
type Exercise struct {
	Index              int    `json:"index"`
	Title              string `json:"title"`
	Notes              string `json:"notes"`
	ExerciseTemplateID string `json:"exercise_template_id"`
	SupersetID         any    `json:"superset_id"`
	Sets               []Set  `json:"sets"`
}
type Workout struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	RoutineID   string     `json:"routine_id"`
	Description string     `json:"description"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     time.Time  `json:"end_time"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CreatedAt   time.Time  `json:"created_at"`
	Exercises   []Exercise `json:"exercises"`
}
type Event struct {
	Type    string  `json:"type"`
	Workout Workout `json:"workout"`
}

func (we *WorkoutEvents) Get(since time.Time) ([]Event, error) {
	var envelope WorkoutEventsEnvelope
	query := url.Values{}
	query.Add("since", since.Format("2006-01-02"))
	path := "/v1/workouts/events?" + query.Encode()
	err := we.client.Call(http.MethodGet, path, nil, &envelope)
	if err != nil {
		return nil, err
	}
	return envelope.Events, nil
}
