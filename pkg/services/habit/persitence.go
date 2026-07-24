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
)

func get(db bob.DB, p core.HabitParams) (habit core.Habit, err error) {
	res, err := buildHabitQuery(p).One(context.Background(), db)
	if err != nil {
		return habit, persistence.CatchDBErr("habits: get", err)
	}
	return bobto.Habit(res), nil
}

func upsert(db bob.DB, p core.Habit) (coreHabit core.Habit, err error) {
	day, err := persistence.GetDay(p.Date, db)
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
	/*
		if f.DayID > 0 {
			q.Apply(models.SelectWhere.Habits.DayID.EQ(int64(f.DayID)))
		}
		if f.CategoryID > 0 {
			q.Apply(models.SelectWhere.Habits.HabitCategoryID.EQ(int64(f.CategoryID)))
		}
	*/
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

/*
func habitCategoryToCore(hc *models.HabitCategory) (c core.HabitCategory) {
	if hc == nil {
		slog.Error("habitCategoryToCore habit category is nil")
		return
	}
	c.ID = uint(hc.ID)
	c.Name = hc.Name
	c.Kind = core.HabitType(hc.Kind)
	c.Description = hc.Description
	// Now we are missing the habits.
	return
}
*/
