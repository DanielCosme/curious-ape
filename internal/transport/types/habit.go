package types

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/sdk/dates"
	"time"
)

type HabitTransport struct {
	ID           int                `json:"id"`
	Success      bool               `json:"success"`
	Date         *time.Time         `json:"date,omitempty"`
	CategoryID   int                `json:"category_id,omitempty"`
	CategoryCode string             `json:"category_code,omitempty"`
	Type         entity.HabitType   `json:"category_type,omitempty"`
	Origin       entity.HabitOrigin `json:"origin,omitempty"`
	Note         string             `json:"note,omitempty"`
	IsAutomated  bool               `json:"is_automated,omitempty"`
}

func (ht *HabitTransport) ToHabit() *entity.Habit {
	return &entity.Habit{
		Success:    ht.Success,
		Origin:     ht.Origin,
		CategoryID: ht.CategoryID,
		Note:       ht.Note,
		Category: &entity.HabitCategory{
			ID:   ht.CategoryID,
			Type: ht.Type,
			Code: ht.CategoryCode,
		},
	}
}

func FromHabitToTransport(h *entity.Habit) *HabitTransport {
	ht := &HabitTransport{
		ID:          h.ID,
		Success:     h.Success,
		Origin:      h.Origin,
		Type:        h.Category.Type,
		IsAutomated: h.IsAutomated,
		Note:        h.Note,
	}

	if h.Day != nil {
		ht.Date = dates.ToPtr(h.Day.Date)
	}

	return ht
}
