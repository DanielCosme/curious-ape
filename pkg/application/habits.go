package application

import (
	"github.com/danielcosme/curious-ape/pkg/core"
	"log/slog"
)

func (a *App) HabitUpsert(params core.UpsertHabitParams) (habit core.Habit, err error) {
	habit, err = a.db.Habits.Upsert(params)
	if err != nil {
		return
	}
	slog.Info("Habit logged",
		"type", habit.Type,
		"state", habit.State,
		"date", habit.Date.Time().Format(core.HumanDateWeekDay),
	)
	return
}
