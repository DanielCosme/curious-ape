package application

import (
	"danicos.dev/daniel/curious-ape/pkg/core"
)

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
