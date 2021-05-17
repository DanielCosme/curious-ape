package data

import "time"

type FoodHabit struct {
	ID    int       `json:"-"` // The - directive hides it from json
	State bool      `json:"state"`
	Date  time.Time `json:"date"`
	Tags  []string  `json:"tags,omitempty"` // Tags: 16/8_fast, lion, calorie_deficit, calorie_surplus.
}
