package day

import (
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/event"
	"danicos.dev/daniel/curious-ape/pkg/oak"
	"danicos.dev/daniel/curious-ape/pkg/persistence"
	"github.com/stephenafamo/bob"
)

type App struct {
	daysDB   core.DayRepository
	habitsDB core.HabitRepository
	bus      event.Bus
}

func New(db bob.DB, event_bus event.Bus) *App {
	return &App{
		daysDB:   persistence.NewDays(db),
		habitsDB: persistence.NewHabits(db),
		bus:      event_bus,
	}
}

func (a *App) Month(date core.Date, order core.OrderParam) ([]core.Day, error) {
	if time.Now().Month() != date.Time.Month() {
		date = date.LastDayOfTheMonth()
	}

	day, err := a.daysDB.Get(core.DayParams{Date: date})
	if core.IfErrNNotFound(err) {
		return nil, err
	}
	daysOfTheMonth := date.RangeMonth()
	if day.IsZero() {
		for _, date := range daysOfTheMonth {
			if _, err := a.GetOrCreate(date); err != nil {
				return nil, err
			}
		}
	}

	return a.daysDB.Find(core.DayParams{Dates: daysOfTheMonth, Order: order})
}

func (a *App) GetOrCreate(d core.Date) (day core.Day, err error) {
	day, err = a.daysDB.Get(core.DayParams{Date: d})
	if core.IfErrNNotFound(err) {
		return
	}
	if day.IsZero() {
		day, err = a.daysDB.Create(d)
		if err != nil {
			return
		}

		hs := []core.Habit{
			{Date: day.Date, State: core.HabitStateNoInfo, Type: core.HabitTypeWakeUp},
			{Date: day.Date, State: core.HabitStateNoInfo, Type: core.HabitTypeFitness},
			{Date: day.Date, State: core.HabitStateNoInfo, Type: core.HabitTypeDeepWork},
			{Date: day.Date, State: core.HabitStateNoInfo, Type: core.HabitTypeEatHealthy},
		}
		for _, habit := range hs {
			err = a.bus.Publish(event.HabitUpsert, habit)
			if err != nil {
				oak.Error("days: error publishing habit upsert", err.Error())
				return day, err
			}
		}
		return a.daysDB.Get(core.DayParams{ID: day.ID})
	}
	return
}
