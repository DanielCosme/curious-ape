package database

import (
	"context"
	"github.com/aarondl/opt/omit"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database/gen/models"
	"github.com/stephenafamo/bob"
)

type Habits struct {
	db bob.DB
}

func (h *Habits) Get(p HabitParams) (core.Habit, error) {
	habit, err := p.BuildQuery(h.db).One()
	if err != nil {
		return core.Habit{}, catchErr("get habit", err)
	}
	return habitToCore(habit), nil
}

func (h *Habits) Create(s models.HabitSetter) (core.Habit, error) {
	habit, err := models.Habits.Insert(context.Background(), h.db, &s)
	if err != nil {
		return core.Habit{}, catchErr("create habit", err)
	}
	ctx := context.Background()
	if err := habit.LoadHabitDay(ctx, h.db); err != nil {
		return core.Habit{}, catchErr("create habit", err)
	}
	if err := habit.LoadHabitHabitCategory(ctx, h.db); err != nil {
		return core.Habit{}, catchErr("create habit", err)
	}
	return habitToCore(habit), nil
}

func (h *Habits) AddLog(s *models.HabitLogSetter) (res core.Habit, err error) {
	_, err = models.HabitLogs.Upsert(
		context.Background(),
		h.db,
		true,
		[]string{"habit_id", "origin"},
		[]string{"success", "detail"},
		s,
	)
	if err != nil {
		return res, err
	}
	res, err = h.Get(HabitParams{ID: s.HabitID.MustGet()})
	if err != nil {
		return res, err
	}
	return res, models.Habits.Update(
		context.Background(),
		h.db,
		&models.HabitSetter{State: omit.From(string(res.State()))},
		&models.Habit{ID: res.ID},
	)
}

func (h *Habits) GetCategory(p HabitCategoryParams) (core.HabitCategory, error) {
	hc, err := p.BuildQuery(h.db).One()
	if err != nil {
		return core.HabitCategory{}, catchErr("get category", err)
	}
	return habitCategoryToCore(hc), nil
}

func habitToCore(m *models.Habit) core.Habit {
	cat := habitCategoryToCore(m.R.HabitCategory)
	logs := habitLogsToCore(m.R.HabitLogs)
	habit := core.NewHabit(core.NewDate(m.R.Day.Date), cat, logs)
	habit.ID = m.ID
	habit.DayID = m.R.Day.ID
	return habit
}

func habitsToCore(ms models.HabitSlice) []core.Habit {
	if len(ms) == 0 {
		return nil
	}
	res := make([]core.Habit, len(ms))
	for idx, h := range ms {
		res[idx] = habitToCore(h)
	}
	return res
}

func habitLogsToCore(ls models.HabitLogSlice) (res []core.HabitLog) {
	for _, l := range ls {
		res = append(res, core.HabitLog{
			ID:          l.ID,
			Success:     l.Success,
			IsAutomated: l.IsAutomated,
			Origin:      core.OriginLog(l.Origin),
			Detail:      l.Detail,
		})
	}
	return
}

func habitCategoryToCore(m *models.HabitCategory) core.HabitCategory {
	return core.HabitCategory{
		ID:          m.ID,
		Name:        m.Name,
		Type:        core.HabitType(m.Type),
		Description: m.Description,
	}
}
