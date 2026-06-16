package persistence

import (
	"context"

	"danicos.dev/daniel/curious-ape/database/gen/dberrors"
	"danicos.dev/daniel/curious-ape/database/gen/models"
	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/oak"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
)

type Habits struct {
	db bob.DB
}

func (h *Habits) Get(p core.HabitParams) (habit core.Habit, err error) {
	res, err := buildHabitQuery(p).One(context.Background(), h.db)
	if err != nil {
		return habit, catchDBErr("habits: get", err)
	}
	return habitToCore(res), nil
}

func (h *Habits) Upsert(p core.Habit) (coreHabit core.Habit, err error) {
	day, err := getDay(p.Date, h.db)
	if err == nil {
		hCategory, err := buildHabitCategoryQuery(core.HabitCategoryParams{Kind: p.Type}).One(context.Background(), h.db)
		if err == nil {
			setter := &models.HabitSetter{
				DayID:           omit.From(day.ID),
				HabitCategoryID: omit.From(hCategory.ID),
				State:           omit.From(string(p.State)),
				NOTE:            omitnull.From(p.Note),
				Automated:       omit.From(p.Automated),
			}
			habit, err := models.Habits.Insert(setter).One(context.Background(), h.db)
			isUpdate := dberrors.HabitErrors.ErrUniqueSqliteAutoindexHabit1.Is(err)
			if err == nil || isUpdate {
				if isUpdate {
					habit, err = models.Habits.Query(
						models.SelectWhere.Habits.DayID.EQ(setter.DayID.GetOrZero()),
						models.SelectWhere.Habits.HabitCategoryID.EQ(setter.HabitCategoryID.GetOrZero()),
						models.Preload.Habit.Day(),
						models.Preload.Habit.HabitCategory(),
					).One(context.Background(), h.db)
					if err == nil {
						// Non-automated habits should not be overwriten by automated ones.
						setterAutomated := setter.Automated.MustGet()
						if habit.Automated == setterAutomated ||
							!setterAutomated ||
							habit.State == string(core.HabitStateNoInfo) {
							if err = habit.Update(context.Background(), h.db, setter); err != nil {
								return coreHabit, catchDBErr("habits: upsert", err)
							}
						} else {
							oak.Info("No-Op UPDATE for habit",
								"current automated", habit.Automated,
								"setter automated", setterAutomated)
						}
					} else {
						return coreHabit, catchDBErr("habits: upsert", err)
					}
				}

				ctx := context.Background()
				if err = habit.LoadDay(ctx, h.db); err == nil {
					if err = habit.LoadHabitCategory(ctx, h.db); err == nil {
						return habitToCore(habit), nil
					}
					return coreHabit, catchDBErr("habits: create: load habit category", err)
				}
				return coreHabit, catchDBErr("habits: create: load habit day", err)
			}
		}
	}
	return coreHabit, catchDBErr("habits: upsert", err)
}

// OLD: 7 err
// OLD: 2 normal
//
// NEW: 5 err
// NEW: 1 normal

func habitToCore(h *models.Habit) (habit core.Habit) {
	if h == nil {
		oak.Error("habitToCore: habit is nil")
		return
	}
	habit.ID = uint(h.ID)
	habit.Date = core.NewDate(h.R.Day.Date)
	habit.State = core.HabitState(h.State)
	habit.Type = core.HabitType(h.R.HabitCategory.Kind)
	habit.Note = h.NOTE.GetOrZero()
	habit.Automated = h.Automated
	return
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
		oak.Error("habitCategoryToCore habit category is nil")
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
