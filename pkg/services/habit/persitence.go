package habit

import (
	"context"
	"log/slog"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/gen/bob/dberrors"
	"danicos.dev/daniel/curious-ape/pkg/gen/bob/models"
	"danicos.dev/daniel/curious-ape/pkg/persistence"
	"danicos.dev/daniel/curious-ape/pkg/persistence/bobto"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
)

func get(db bob.DB, p core.HabitParams) (habit core.Habit, err error) {
	res, err := buildHabitQuery(p).One(context.Background(), db)
	if err != nil {
		return habit, persistence.CatchDBErr("habits: get", err)
	}
	return bobto.Habit(res), nil
}

func find(db bob.DB, p core.HabitParams) ([]core.Habit, error) {
	res, err := buildHabitQuery(p).All(context.Background(), db)
	if err != nil {
		return nil, persistence.CatchDBErr("habits: find", err)
	}
	habits := make([]core.Habit, 0, len(res))
	for _, h := range res {
		habits = append(habits, bobto.Habit(h))
	}
	return habits, nil
}

func upsert(db bob.DB, p core.Habit) (coreHabit core.Habit, err error) {
	day, err := persistence.GetDay(db, p.Date)
	if err == nil {
		hCategory, err := buildHabitCategoryQuery(core.HabitCategoryParams{Kind: p.Type}).One(context.Background(), db)
		if err == nil {
			setter := &models.HabitSetter{
				DayID:           omit.From(day.ID),
				HabitCategoryID: omit.From(hCategory.ID),
				State:           omit.From(string(p.State)),
				NOTE:            omitnull.From(p.Note),
				Automated:       omit.From(p.Automated),
			}
			habit, err := models.Habits.Insert(setter).One(context.Background(), db)
			isUpdate := dberrors.HabitErrors.ErrUniqueSqliteAutoindexHabit1.Is(err)
			if err == nil || isUpdate {
				if isUpdate {
					habit, err = models.Habits.Query(
						models.SelectWhere.Habits.DayID.EQ(setter.DayID.GetOrZero()),
						models.SelectWhere.Habits.HabitCategoryID.EQ(setter.HabitCategoryID.GetOrZero()),
						models.Preload.Habit.Day(),
						models.Preload.Habit.HabitCategory(),
					).One(context.Background(), db)
					if err == nil {
						// Non-automated habits should not be overwriten by automated ones.
						setterAutomated := setter.Automated.MustGet()
						if habit.Automated == setterAutomated ||
							!setterAutomated ||
							habit.State == string(core.HabitStateNoInfo) {
							if err = habit.Update(context.Background(), db, setter); err != nil {
								return coreHabit, persistence.CatchDBErr("habits: upsert", err)
							}
						} else {
							slog.Info("No-Op UPDATE for habit",
								"current automated", habit.Automated,
								"setter automated", setterAutomated)
						}
					} else {
						return coreHabit, persistence.CatchDBErr("habits: upsert", err)
					}
				}

				ctx := context.Background()
				if err = habit.LoadDay(ctx, db); err == nil {
					if err = habit.LoadHabitCategory(ctx, db); err == nil {
						return bobto.Habit(habit), nil
					}
					return coreHabit, persistence.CatchDBErr("habits: create: load habit category", err)
				}
				return coreHabit, persistence.CatchDBErr("habits: create: load habit day", err)
			}
		}
	}
	return coreHabit, persistence.CatchDBErr("habits: upsert", err)
}

func buildHabitQuery(f core.HabitParams) *sqlite.ViewQuery[*models.Habit, models.HabitSlice] {
	q := models.Habits.Query()
	q.Apply(models.Preload.Habit.Day())
	q.Apply(models.Preload.Habit.HabitCategory())
	if f.ID > 0 {
		q.Apply(models.SelectWhere.Habits.ID.EQ(int64(f.ID)))
	}
	// maps to habit_category.kind (e.g. wake_up, fitness).
	if f.Type != "" {
		q.Apply(models.SelectJoins().Habits.InnerJoin.HabitCategory)
		q.Apply(models.SelectWhere.HabitCategories.Kind.EQ(string(f.Type)))
	}
	if !f.From.Time.IsZero() && !f.To.Time.IsZero() {
		q.Apply(models.SelectJoins().Habits.InnerJoin.Day)
		q.Apply(models.SelectWhere.Days.Date.GTE(f.From.Time))
		q.Apply(models.SelectWhere.Days.Date.LTE(f.To.Time))
		q.Apply(sm.OrderBy(models.Days.Columns.Date).Desc())
	}
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
