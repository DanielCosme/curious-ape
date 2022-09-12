package types

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/go-sdk/dates"
	"time"
)

type HabitTransport struct {
	ID           int                  `json:"id"`
	Date         *time.Time           `json:"date,omitempty"`
	Status       entity.HabitStatus   `json:"status"`
	Success      bool                 `json:"success,omitempty"`
	CategoryID   int                  `json:"category_id,omitempty"`
	CategoryCode string               `json:"category_code,omitempty"`
	Type         entity.HabitType     `json:"category_type,omitempty"`
	Origin       entity.DataSource    `json:"origin,omitempty"`
	Note         string               `json:"note,omitempty"`
	IsAutomated  bool                 `json:"is_automated,omitempty"`
	Logs         []*HabitTransportLog `json:"logs,omitempty"`
}

type HabitTransportLog struct {
	Success     bool              `json:"success"`
	Origin      entity.DataSource `json:"origin"`
	Note        string            `json:"note,omitempty"`
	IsAutomated bool              `json:"is_automated"`
}

func (ht *HabitTransport) ToHabit() *entity.Habit {
	return &entity.Habit{
		CategoryID: ht.CategoryID,
		Logs: []*entity.HabitLog{
			{
				Success:     ht.Success,
				Origin:      ht.Origin,
				Note:        ht.Note,
				IsAutomated: ht.IsAutomated,
			},
		},
	}
}

func FromHabitToTransport(h *entity.Habit) *HabitTransport {
	if h != nil {
		ht := &HabitTransport{
			ID:     h.ID,
			Type:   h.Category.Type,
			Status: h.Status,
		}

		if h.Day != nil {
			ht.Date = dates.Ptr(h.Day.Date)
		}

		for _, l := range h.Logs {
			ht.Logs = append(ht.Logs, fromHabitLogToTransport(l))
		}

		return ht
	}
	return nil
}

func FromHabitToTransportSlice(hs []*entity.Habit) []*HabitTransport {
	hts := []*HabitTransport{}
	for _, h := range hs {
		hts = append(hts, FromHabitToTransport(h))
	}
	return hts
}

func fromHabitLogToTransport(hl *entity.HabitLog) *HabitTransportLog {
	return &HabitTransportLog{
		Success:     hl.Success,
		Origin:      hl.Origin,
		Note:        hl.Note,
		IsAutomated: hl.IsAutomated,
	}
}
