package application

import (
	"github.com/aarondl/opt/omit"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/database"
	"github.com/danielcosme/curious-ape/pkg/database/gen/models"
	"log/slog"
)

func (a *App) HabitUpsert(date core.Date, hk core.HabitKind, state core.HabitState) (*models.Habit, error) {
	hc, err := a.db.Habits.GetCategory(database.HabitCategoryParams{Kind: hk})
	if err != nil {
		return nil, err
	}
	day, err := a.db.Days.GetOrCreate(database.DayParams{Date: date})
	if err != nil {
		return nil, err
	}
	habit, err := a.db.Habits.Upsert(&models.HabitSetter{
		DayID:           omit.From(day.ID),
		HabitCategoryID: omit.From(hc.ID),
		State:           omit.From(string(state)),
	})
	if err != nil {
		return nil, err
	}
	slog.Info("Habit logged",
		"name", hc.Name,
		"state", habit.State,
		"date", day.Date.Format(core.HumanDateWeekDay))
	return habit, nil
}
