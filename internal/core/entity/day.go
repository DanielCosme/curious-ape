package entity

import "time"

// Day is the top level entity of the domain
type Day struct {
	Entity
	Date time.Time
}
