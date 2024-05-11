package database

import (
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database/gen/models"
)

func habitToCore(m *models.Habit) core.Habit {
	h := core.Habit{}
	return h
}

func habitsToCore(ms models.HabitSlice) []core.Habit {
	res := make([]core.Habit, len(ms))
	for idx, d := range ms {
		res[idx] = habitToCore(d)
	}
	return res
}
