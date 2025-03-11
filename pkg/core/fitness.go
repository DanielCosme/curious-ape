package core

import (
	"fmt"
	"time"
)

type FitnessLog struct {
	ID    int32
	DayID int32

	Title     string
	Type      string
	Date      Date
	StartTime time.Time
	EndTime   time.Time
	Origin    Integration
	Raw       string
	Note      string
	// Generated
	Day *Day
}

func NewFitnessLog(day Day) FitnessLog {
	ft := FitnessLog{
		DayID: day.ID,
		Date:  day.Date,
		Day:   &day,
	}
	return ft
}

func (d *FitnessLog) ToHabitLogFitness() (res NewHabitParams) {
	res.Date = d.Date
	res.HabitType = HabitTypeExercise
	res.Origin = OriginLogFitness
	res.Automated = true
	dur := d.EndTime.Sub(d.StartTime).Round(time.Minute)
	res.Detail = fmt.Sprintf("%s - %s (%s)", d.StartTime.Format(Time), d.EndTime.Format(Time), dur)
	res.Success = true
	return
}
