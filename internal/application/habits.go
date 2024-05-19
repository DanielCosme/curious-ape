package application

import (
	"errors"
	"github.com/aarondl/opt/omit"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database/gen/models"
	"log/slog"
)

func (a *App) HabitUpsert(params core.NewHabitParams) (habit core.Habit, err error) {
	if !params.Valid() {
		return habit, errors.New("invalid habit params")
	}
	habit, err = a.db.HabitGetOrCreate(params.Date, params.HabitType)
	if err != nil {
		return habit, err
	}
	habit, err = a.db.Habits.AddLog(&models.HabitLogSetter{
		HabitID:     omit.From(habit.ID),
		Origin:      omit.From(string(params.Origin)),
		Success:     omit.From(params.Success),
		IsAutomated: omit.From(params.Automated),
		Detail:      omit.From(params.Detail),
	})
	slog.Info("Habit logged",
		"name", habit.Category.Name,
		"state", habit.State(),
		"origin", params.Origin,
		"detail", params.Detail,
	)
	return
}
