package application

import (
	"context"

	"git.danicos.dev/daniel/curious-ape/pkg/core"
	"git.danicos.dev/daniel/curious-ape/pkg/oak"
)

func (a *App) HabitUpsert(ctx context.Context, params core.Habit) (habit core.Habit, err error) {
	logger := oak.FromContext(ctx)

	habit, err = a.db.Habits.Upsert(params)
	if err != nil {
		return
	}
	logger.Info("Habit logged",
		"type", habit.Type,
		"state", habit.State,
		"date", habit.Date.Time().Format(core.HumanDateWeekDay),
	)
	return
}

func (a *App) HabitFlip(id int) (habit core.Habit, err error) {
	habit, err = a.db.Habits.Get(core.HabitParams{ID: id})
	if err != nil {
		return
	}
	state := core.HabitStateNotDone
	if habit.State == core.HabitStateNotDone || habit.State == core.HabitStateNoInfo {
		state = core.HabitStateDone
	}
	habit.State = state
	return a.db.Habits.Upsert(core.Habit{
		Date:  habit.Date,
		Type:  habit.Type,
		State: habit.State,
	})
}
