package core

import "time"

type DeepWorkLog struct {
	ID          int32
	DayID       int32
	Date        Date
	Duration    time.Duration
	Origin      Integration
	IsAutomated bool
}

func NewDeepWorkLog(dur time.Duration, day Day) DeepWorkLog {
	return DeepWorkLog{
		Duration: dur,
		DayID:    day.ID,
		Date:     day.Date,
	}
}

func (d *DeepWorkLog) IsZero() bool {
	return d.Duration.Seconds() == 0
}

func (d *DeepWorkLog) ToHabitLogDeepWork() (res NewHabitParams) {
	res.Date = d.Date
	res.HabitType = HabitTypeDeepWork
	res.Origin = OriginLogDeepWork
	res.Automated = true
	res.Detail = d.Duration.String()
	if d.Duration > time.Hour*5 {
		res.Success = true
	}
	return
}
