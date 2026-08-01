package habit

import (
	"errors"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/event"
	"danicos.dev/daniel/curious-ape/pkg/oak"
	"danicos.dev/daniel/curious-ape/pkg/persistence"
	"github.com/stephenafamo/bob"
)

type App struct {
	habitsDB core.HabitRepository
	bus      event.Bus
}

func New(db bob.DB, bus event.Bus) *App {
	app := &App{
		habitsDB: persistence.NewHabits(db),
		bus:      bus,
	}
	bus.Subscribe(event.HabitUpsert, app.EventHandleHabitUpsert)
	return app
}

func (a *App) EventHandleHabitUpsert(data any) error {
	params, ok := (data).(core.Habit)
	if !ok {
		e := "invalid event data type"
		oak.Error(e)
		return errors.New(e)
	}
	_, err := a.HabitUpsert(params)
	return err
}

func (a *App) HabitUpsert(params core.Habit) (habit core.Habit, err error) {
	habit, err = a.habitsDB.Upsert(params)
	if err == nil && habit.State != core.HabitStateNoInfo {
		oak.Info("Habit logged",
			"type", habit.Type,
			"state", habit.State,
			"date", habit.Date.Time.Format(core.HumanDateWeekDay),
		)
	}
	return
}
