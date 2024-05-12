package application

import (
	"errors"
	"github.com/aarondl/opt/omit"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database/gen/models"
	"log/slog"
)

func (a *App) HabitUpsert(params core.HabitParams) (habit core.Habit, err error) {
	if !params.Valid() {
		return habit, errors.New("invalid habit params")
	}
	habit, err = a.db.HabitGetOrCreate(params.Date, params.CategoryID)
	if err != nil {
		return habit, err
	}
	habit, err = a.db.Habits.AddLog(&models.HabitLogSetter{
		HabitID:     omit.From(habit.ID),
		Origin:      omit.From(string(params.Origin)),
		Success:     omit.From(params.Success),
		IsAutomated: omit.From(params.Automated),
	})
	slog.Info("Habit logged", "name", habit.Category.Name, "state", habit.State())
	return
}
