package core

import "time"

type DeepWorkLog struct {
	Day      Date
	Duration time.Duration
}

func NewDeepWorkLog(d time.Duration, day Date) DeepWorkLog {
	return DeepWorkLog{
		Day:      day,
		Duration: d,
	}
}

func (d *DeepWorkLog) IsZero() bool {
	return d.Duration.Seconds() == 0
}

func (d *DeepWorkLog) ToHabitLogDeepWork() (res NewHabitParams) {
	res.Date = d.Day
	res.HabitType = HabitTypeDeepWork
	res.Origin = OriginLogDeepWork
	res.Automated = true
	res.Detail = d.Duration.String()
	if d.Duration > time.Hour*5 {
		res.Success = true
	}
	return
}
