package database

import (
	"context"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database/gen/models"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
)

type DayParams struct {
	ID    int32
	Date  core.Date
	Dates core.DateSlice
	R     []Relation
}

func (f DayParams) BuildQuery(exec bob.Executor) *sqlite.ViewQuery[*models.Day, models.DaySlice] {
	q := models.Days.Query(context.Background(), exec)
	if f.ID > 0 {
		q.Apply(models.SelectWhere.Days.ID.EQ(f.ID))
	}
	if !f.Date.Time().IsZero() {
		q.Apply(models.SelectWhere.Days.Date.EQ(f.Date.Time()))
	}
	if len(f.Dates) > 0 {
		q.Apply(models.SelectWhere.Days.Date.In(f.Dates.ToTimeSlice()...))
	}
	for _, r := range f.R {
		switch r {
		case RelationHabit:
			q.Apply(models.ThenLoadDayHabits())
		case RelationSleep:
			q.Apply(models.ThenLoadDaySleepLogs())
		case RelationFitness:
			q.Apply(models.ThenLoadDayFitnessLogs())
		}
	}
	return q
}

type HabitParams struct {
	ID         int32
	DayID      int32
	CategoryID int32
}

func (f HabitParams) BuildQuery(exec bob.Executor) *sqlite.ViewQuery[*models.Habit, models.HabitSlice] {
	q := models.Habits.Query(context.Background(), exec)
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
	q.Apply(models.ThenLoadHabitHabitLogs())
	q.Apply(models.PreloadHabitHabitCategory())
	return q
}

type HabitCategoryParams struct {
	ID int32
}

func (f HabitCategoryParams) BuildQuery(exec bob.Executor) *sqlite.ViewQuery[*models.HabitCategory, models.HabitCategorySlice] {
	q := models.HabitCategories.Query(context.Background(), exec)
	if f.ID > 0 {
		q.Apply(models.SelectWhere.HabitCategories.ID.EQ(f.ID))
	}
	return q
}
