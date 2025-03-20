package database

import (
	"context"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/database/gen/models"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
)

type Habits struct {
	db bob.DB
}

func (h *Habits) Get(p HabitParams) (*models.Habit, error) {
	habit, err := p.BuildQuery().One(context.Background(), h.db)
	if err != nil {
		return nil, catchDBErr("habits: get", err)
	}
	return habit, nil
}

func (h *Habits) Upsert(s *models.HabitSetter) (*models.Habit, error) {
	habit, err := models.Habits.Insert(s).One(context.Background(), h.db)
	if err != nil {
		if models.HabitErrors.ErrUniqueDayIdAndHabitCategoryId.Is(err) {
			habit, err = models.Habits.
				Update(s.UpdateMod(),
					models.UpdateWhere.Habits.DayID.EQ(s.DayID.GetOrZero()),
					models.UpdateWhere.Habits.HabitCategoryID.EQ(s.HabitCategoryID.GetOrZero()),
				).
				One(context.Background(), h.db)
			if err != nil {
				return nil, catchDBErr("habits: create", err)
			}
		} else {
			return nil, catchDBErr("habits: create", err)
		}
	}
	ctx := context.Background()
	if err := habit.LoadHabitDay(ctx, h.db); err != nil {
		return nil, catchDBErr("habits: create: load habit day", err)
	}
	if err := habit.LoadHabitHabitCategory(ctx, h.db); err != nil {
		return nil, catchDBErr("habits: create: load habit category", err)
	}
	return habit, nil
}

func (h *Habits) GetCategory(p HabitCategoryParams) (*models.HabitCategory, error) {
	hc, err := p.BuildQuery().One(context.Background(), h.db)
	if err != nil {
		return nil, catchDBErr("habit_category: get", err)
	}
	return hc, nil
}

type HabitParams struct {
	ID         int32
	DayID      int32
	CategoryID int32
}

func (f HabitParams) BuildQuery() *sqlite.ViewQuery[*models.Habit, models.HabitSlice] {
	q := models.Habits.Query()
	if f.ID > 0 {
		q.Apply(models.SelectWhere.Habits.ID.EQ(f.ID))
	}
	if f.DayID > 0 {
		q.Apply(models.SelectWhere.Habits.DayID.EQ(f.DayID))
	}
	if f.CategoryID > 0 {
		q.Apply(models.SelectWhere.Habits.HabitCategoryID.EQ(f.CategoryID))
	}
	q.Apply(models.PreloadHabitDay())
	q.Apply(models.PreloadHabitHabitCategory())
	return q
}

type HabitCategoryParams struct {
	ID   int32
	Kind core.HabitType
}

func (f HabitCategoryParams) BuildQuery() *sqlite.ViewQuery[*models.HabitCategory, models.HabitCategorySlice] {
	q := models.HabitCategories.Query()
	if f.ID > 0 {
		q.Apply(models.SelectWhere.HabitCategories.ID.EQ(f.ID))
	}
	if f.Kind != "" {
		q.Apply(models.SelectWhere.HabitCategories.Kind.EQ(string(f.Kind)))
	}
	return q
}
