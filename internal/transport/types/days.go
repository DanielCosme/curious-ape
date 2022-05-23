package types

import "github.com/danielcosme/curious-ape/internal/core/entity"

type DayTransport struct {
	Date   string            `json:"date"`
	Habits []*HabitTransport `json:"habits,omitempty"`
}

func DayToTransport(d *entity.Day) *DayTransport {
	dt := &DayTransport{
		Date:   entity.FormatDate(d.Date),
	}

	if len(d.Habits) > 0 {
		for _, h := range d.Habits {
			dt.Habits = append(dt.Habits, FromHabitToTransport(h))
		}
	}

	return dt
}
