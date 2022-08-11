package types

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"time"
)

type DayTransport struct {
	Date            string                 `json:"date"`
	DeepWork        string                 `json:"deep_work"`
	DeepWorkMinutes int                    `json:"deep_work_minutes,omitempty"`
	Habits          []*HabitTransport      `json:"habits,omitempty"`
	SleepLogs       []*SleepLogTransport   `json:"sleep_logs,omitempty"`
	FitnessLogs     []*FitnessLogTransport `json:"fitnessLogs,omitempty"`
}

func (dt *DayTransport) ToDay() *entity.Day {
	return &entity.Day{
		DeepWorkMinutes: dt.DeepWorkMinutes,
	}
}

func DayToTransport(d *entity.Day) *DayTransport {
	dt := &DayTransport{
		Date:     entity.FormatDate(d.Date),
		DeepWork: (time.Duration(d.DeepWorkMinutes) * time.Minute).String(),
	}

	if len(d.Habits) > 0 {
		for _, h := range d.Habits {
			dt.Habits = append(dt.Habits, FromHabitToTransport(h))
		}
	}

	if len(d.SleepLogs) > 0 {
		dt.SleepLogs = FromSleepLogToTransportSlice(d.SleepLogs)
	}

	if len(d.FitnessLogs) > 0 {
		dt.FitnessLogs = FromFitnessLogToTransportSlice(d.FitnessLogs)
	}

	return dt
}
