package persistence

import (
	"context"
	"log/slog"

	"github.com/aarondl/opt/omit"
	"github.com/danielcosme/curious-ape/database/gen/dberrors"
	"github.com/danielcosme/curious-ape/database/gen/models"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
)

type Habits struct {
	db bob.DB
}

func (h Habits) Get(p core.HabitParams) (habit core.Habit, err error) {
	res, err := buildHabitQuery(p).One(context.Background(), h.db)
	if err != nil {
		return habit, catchDBErr("habits: get", err)
	}
	return habitToCore(res), nil
}

func habitToCore(h *models.Habit) (habit core.Habit) {
	if h == nil {
		slog.Error("habitToCore habit is nil")
		return
	}
	habit.ID = int(h.ID)
	habit.Date = core.NewDate(h.R.Day.Date)
	habit.State = core.HabitState(h.State)
	habit.Type = core.HabitType(h.R.HabitCategory.Kind)
	habit.Automated = h.Automated
	return
}

func habitCategoryToCore(hc *models.HabitCategory) (c core.HabitCategory) {
	if hc == nil {
		slog.Error("habitCategoryToCore habit category is nil")
		return
	}
	c.ID = int(hc.ID)
	c.Name = hc.Name
	c.Kind = core.HabitType(hc.Kind)
	c.Description = hc.Description
	return
}

func (h Habits) Upsert(p core.UpsertHabitParams) (coreHabit core.Habit, err error) {
	day, err := BuildDayQuery(core.DayParams{Date: p.Date, Lite: true}).One(context.Background(), h.db)
	if err != nil {
		return coreHabit, catchDBErr("habits: upsert", err)
	}
	hCategory, err := buildHabitCategoryQuery(core.HabitCategoryParams{Kind: p.Type}).One(context.Background(), h.db)
	if err != nil {
		return coreHabit, catchDBErr("habits: upsert", err)
	}
	s := &models.HabitSetter{
		DayID:           omit.From(day.ID),
		HabitCategoryID: omit.From(hCategory.ID),
		State:           omit.From(string(p.State)),
		Automated:       omit.From(p.Automated),
	}
	habit, err := models.Habits.Insert(s).One(context.Background(), h.db)
	if err != nil {
		if dberrors.HabitErrors.ErrUniqueSqliteAutoindexHabit1.Is(err) {
			habit, err = models.Habits.Query(
				models.SelectWhere.Habits.DayID.EQ(s.DayID.GetOrZero()),
				models.SelectWhere.Habits.HabitCategoryID.EQ(s.HabitCategoryID.GetOrZero()),
			).One(context.Background(), h.db)
			if err != nil {
				return coreHabit, catchDBErr("habits: upsert", err)
			}

			// No-op for a non-automated habit for which the update is automated.
			if !habit.Automated && s.Automated.GetOrZero() {
				return habitToCore(habit), nil
			}
			err = habit.Update(context.Background(), h.db, s)
			if err != nil {
				return coreHabit, catchDBErr("habits: upsert", err)
			}
		} else {
			return coreHabit, catchDBErr("habits: create", err)
		}
	}
	ctx := context.Background()
	if err := habit.LoadDay(ctx, h.db); err != nil {
		return coreHabit, catchDBErr("habits: create: load habit day", err)
	}
	if err := habit.LoadHabitCategory(ctx, h.db); err != nil {
		return coreHabit, catchDBErr("habits: create: load habit category", err)
	}
	return habitToCore(habit), nil
}

func buildHabitQuery(f core.HabitParams) *sqlite.ViewQuery[*models.Habit, models.HabitSlice] {
	q := models.Habits.Query()
	if f.ID > 0 {
		q.Apply(models.SelectWhere.Habits.ID.EQ(int64(f.ID)))
	}
	if f.DayID > 0 {
		q.Apply(models.SelectWhere.Habits.DayID.EQ(int64(f.DayID)))
	}
	if f.CategoryID > 0 {
		q.Apply(models.SelectWhere.Habits.HabitCategoryID.EQ(int64(f.CategoryID)))
	}
	q.Apply(models.Preload.Habit.Day())
	q.Apply(models.Preload.Habit.HabitCategory())
	return q
}

func buildHabitCategoryQuery(f core.HabitCategoryParams) *sqlite.ViewQuery[*models.HabitCategory, models.HabitCategorySlice] {
	q := models.HabitCategories.Query()
	if f.ID > 0 {
		q.Apply(models.SelectWhere.HabitCategories.ID.EQ(int64(f.ID)))
	}
	if f.Kind != "" {
		q.Apply(models.SelectWhere.HabitCategories.Kind.EQ(string(f.Kind)))
	}
	return q
}
