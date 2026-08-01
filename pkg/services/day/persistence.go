package day

import (
	"context"
	"fmt"
	"slices"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/gen/bob/models"
	"danicos.dev/daniel/curious-ape/pkg/persistence"
	"danicos.dev/daniel/curious-ape/pkg/persistence/bobto"
	"github.com/aarondl/opt/omit"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
)

func Find(db bob.DB, p core.DayParams) (days []core.Day, err error) {
	res, err := BuildDayQuery(p).All(context.Background(), db)
	if err == nil {

		for _, day := range res { // TODO: optimize this.
			if slices.Contains(p.WithRelation, core.DayRelationHabits) {
				if err = loadHabitRelations(db, day); err != nil {
					return days, persistence.CatchDBErr("days: find", err)
				}
			}

			days = append(days, bobto.Day(day))
		}
		return
	} else {
		return days, persistence.CatchDBErr("days: find", err)
	}
}

func new(db bob.DB, date core.Date) (day core.Day, err error) {
	s := &models.DaySetter{Date: omit.From(date.Time)}
	res, err := models.Days.Insert(s).One(context.Background(), db)
	return bobto.Day(res), err
}

func get(db bob.DB, p core.DayParams) (day core.Day, err error) {
	res, err := BuildDayQuery(p).One(context.Background(), db)
	if err == nil {
		if slices.Contains(p.WithRelation, core.DayRelationHabits) {
			err = loadHabitRelations(db, res)
			if err != nil {
				return day, persistence.CatchDBErr("days: get", err)
			}
		}
		return bobto.Day(res), err
	}
	return day, persistence.CatchDBErr("days: get", err)
}

func loadHabitRelations(db bob.DB, m *models.Day) (err error) {
	if err = m.R.Habits.LoadDay(context.Background(), db); err == nil {
		if err = m.R.Habits.LoadHabitCategory(context.Background(), db); err == nil {
			return nil
		}
	}
	return persistence.CatchDBErr("days: load: habit relations", err)
}

func BuildDayQuery(f core.DayParams) *sqlite.ViewQuery[*models.Day, models.DaySlice] {
	q := models.Days.Query()

	if f.ID > 0 {
		q.Apply(models.SelectWhere.Days.ID.EQ(int64(f.ID)))
	}

	if !f.Date.Time.IsZero() {
		q.Apply(models.SelectWhere.Days.Date.EQ(f.Date.Time))
	}

	if len(f.Dates) > 0 {
		q.Apply(models.SelectWhere.Days.Date.In(f.Dates.ToTimeSlice()...))
	}

	for _, relation := range f.WithRelation {
		switch relation {
		case core.DayRelationHabits:
			q.Apply(models.SelectThenLoad.Day.Habits())
		case core.DayRelationSleepLogs:
			q.Apply(models.SelectThenLoad.Day.SleepLogs())
		case core.DayRelationFitnessLogs:
			q.Apply(models.SelectThenLoad.Day.FitnessLogs())
		case core.DayRelationDeepWorkLogs:
			q.Apply(
				models.SelectThenLoad.Day.DeepWorkLogs(
					sm.OrderBy(models.DeepWorkLogs.Columns.StartTime).Desc(),
				),
			)
		default:
			panic(fmt.Sprintf("unexpected core.DayRelations: %#v", relation))
		}
	}

	if f.Order == core.DESC {
		q.Apply(sm.OrderBy(models.Days.Columns.Date).Desc())
	}

	return q
}
