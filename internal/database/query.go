package database

import (
	"context"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database/gen/models"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
)

type DayParams struct {
	Date  core.Date
	Dates core.DateSlice
	R     []Relation
}

func (f DayParams) BuildQuery(exec bob.Executor) *sqlite.ViewQuery[*models.Day, models.DaySlice] {
	q := models.Days.Query(context.Background(), exec)

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
