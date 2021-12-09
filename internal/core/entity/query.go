package entity

import "time"

type Query struct {
	Offset int
	Limit  int
}

type DateQuery struct {
	From time.Time
	To   time.Time
}
