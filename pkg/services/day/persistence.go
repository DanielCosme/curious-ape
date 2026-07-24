package day

import (
	"context"
	"log/slog"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/gen/bob/models"
	"danicos.dev/daniel/curious-ape/pkg/persistence"
	"github.com/aarondl/opt/omit"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
)

func new(db bob.DB, date core.Date) (day core.Day, err error) {
	s := &models.DaySetter{Date: omit.From(date.Time())}
	res, err := models.Days.Insert(s).One(context.Background(), db)
	return dayToCore(res), err
}

func get(db bob.DB, p core.DayParams) (day core.Day, err error) {
	res, err := BuildDayQuery(p).One(context.Background(), db)
	if err == nil {
		err = loadHabitRelations(db, res)
		return dayToCore(res), err
	}
	return day, persistence.CatchDBErr("days: get", err)
}

func getOrCreate(db bob.DB, params core.DayParams) (day core.Day, err error) {
	day, err = get(db, params)
	if core.IfErrNNotFound(err) {
		return
	}
	if day.IsZero() {
		day, err = new(db, params.Date)
	}
	return
}

func find(db bob.DB, p core.DayParams) (days []core.Day, err error) {
	res, err := BuildDayQuery(p).All(context.Background(), db)
	if err == nil {
		for _, day := range res { // TODO: optimize this.
			if err = loadHabitRelations(db, day); err == nil {
				days = append(days, dayToCore(day))
			} else {
				return days, persistence.CatchDBErr("days: find", err)
			}
		}
		return
	} else {
		return days, persistence.CatchDBErr("days: find", err)
	}
}

func loadHabitRelations(db bob.DB, m *models.Day) (err error) {
	if err = m.R.Habits.LoadDay(context.Background(), db); err == nil {
		if err = m.R.Habits.LoadHabitCategory(context.Background(), db); err == nil {
			return nil
		}
	}
	return persistence.CatchDBErr("days: load: habit relations", err)
}

func dayToCore(d *models.Day) (day core.Day) {
	if d == nil {
		slog.Error("dayToCore: day is nil")
		return
	}
	day.ID = uint(d.ID)
	day.Date = core.NewDate(d.Date)
	// for _, h := range d.R.Habits {
	// 	habit := habitToCore(h)
	// 	switch habit.Type {
	// 	case core.HabitTypeWakeUp:
	// 		day.Habits.Sleep = habit
	// 	case core.HabitTypeFitness:
	// 		day.Habits.Fitness = habit
	// 	case core.HabitTypeDeepWork:
	// 		day.Habits.DeepWork = habit
	// 	case core.HabitTypeEatHealthy:
	// 		day.Habits.Eat = habit
	// 	}
	// 	if h.State == string(core.HabitStateDone) {
	// 		day.Habits.Score += 1
	// 	}
	// 	day.Habits.Hs = append(day.Habits.Hs, habit)
	// }
	// for _, sl := range d.R.SleepLogs {
	// 	day.SleepLogs = append(day.SleepLogs, sleepLogToCore(d, sl))
	// }
	// for _, fl := range d.R.FitnessLogs {
	// 	day.FitnessLogs = append(day.FitnessLogs, fitnessLogToCore(d, fl))
	// }
	// for _, wl := range d.R.DeepWorkLogs {
	// 	day.DeepWorkLogs = append(day.DeepWorkLogs, deepWorkLogToCore(d, wl))
	// }
	return day
}

func BuildDayQuery(f core.DayParams) *sqlite.ViewQuery[*models.Day, models.DaySlice] {
	q := models.Days.Query()
	if f.ID > 0 {
		q.Apply(models.SelectWhere.Days.ID.EQ(int64(f.ID)))
	}
	if !f.Date.Time().IsZero() {
		q.Apply(models.SelectWhere.Days.Date.EQ(f.Date.Time()))
	}
	if len(f.Dates) > 0 {
		q.Apply(models.SelectWhere.Days.Date.In(f.Dates.ToTimeSlice()...))
	}
	q.Apply(models.SelectThenLoad.Day.Habits())
	q.Apply(models.SelectThenLoad.Day.SleepLogs())
	q.Apply(models.SelectThenLoad.Day.FitnessLogs())
	q.Apply(
		models.SelectThenLoad.Day.DeepWorkLogs(
			sm.OrderBy(models.DeepWorkLogs.Columns.StartTime).Desc(),
		),
	)
	if f.Order == core.DESC {
		q.Apply(sm.OrderBy(models.Days.Columns.Date).Desc())
	}
	return q
}

func getDay(date core.Date, exec bob.Executor) (*models.Day, error) {
	return BuildDayQuery(core.DayParams{Date: date}).One(context.Background(), exec)
}
