package types

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"time"
)

type HabitTransport struct {
	Success      bool               `json:"success"`
	Date         time.Time          `json:"date,omitempty"`
	CategoryID   int                `json:"category_id,omitempty"`
	CategoryCode string             `json:"category_code,omitempty"`
	Type         entity.HabitType   `json:"category_type,omitempty"`
	Origin       entity.HabitOrigin `json:"origin,omitempty"`
	Note         string             `json:"note,omitempty"`
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
	return &HabitTransport{
		Success: h.Success,
		Origin:  h.Origin,
		Date:    h.Day.Date,
		Type:    h.Category.Type,
	}
}
