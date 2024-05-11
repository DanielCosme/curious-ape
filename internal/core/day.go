package core

type Day struct {
	ID   int32
	Date Date
	// generated
	Habits []Habit
}

func (d *Day) IsZero() bool {
	return d.Date.Time().IsZero()
}
