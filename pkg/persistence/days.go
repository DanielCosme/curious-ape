package persistence

import (
	"context"
	"log/slog"

	"danicos.dev/daniel/curious-ape/database/gen/models"
	"danicos.dev/daniel/curious-ape/pkg/core"
	"github.com/aarondl/opt/omit"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
)

type Days struct {
	db bob.DB
}

func NewDays(executor bob.DB) *Days {
	return &Days{db: executor}
}

func (d *Days) Create(date core.Date) (day core.Day, err error) {
	s := &models.DaySetter{Date: omit.From(date.Time())}
	res, err := models.Days.Insert(s).One(context.Background(), d.db)
	return dayToCore(res), err
}

func (d *Days) Get(p core.DayParams) (day core.Day, err error) {
	res, err := BuildDayQuery(p).One(context.Background(), d.db)
	if err == nil {
		err = d.LoadHabitRelations(res)
		return dayToCore(res), err
	}
	return day, catchDBErr("days: get", err)
}

func (d *Days) GetOrCreate(p core.DayParams) (day core.Day, err error) {
	day, err = d.Get(p)
	if core.IfErrNNotFound(err) {
		day, err = d.Create(p.Date)
	}
	return
}

func (d *Days) Find(p core.DayParams) (days []core.Day, err error) {
	res, err := BuildDayQuery(p).All(context.Background(), d.db)
	if err == nil {
		for _, day := range res { // TODO: optimize this.
			if err = d.LoadHabitRelations(day); err == nil {
				days = append(days, dayToCore(day))
			} else {
				return days, catchDBErr("days: find", err)
			}
		}
		return
	} else {
		return days, catchDBErr("days: find", err)
	}
}

func (d *Days) LoadHabitRelations(m *models.Day) (err error) {
	if err = m.R.Habits.LoadDay(context.Background(), d.db); err == nil {
		if err = m.R.Habits.LoadHabitCategory(context.Background(), d.db); err == nil {
			return nil
		}
	}
	return catchDBErr("days: load: habit relations", err)
}

func dayToCore(d *models.Day) (day core.Day) {
	if d == nil {
		slog.Error("dayToCore: day is nil")
		return
	}
	day.ID = uint(d.ID)
	day.Date = core.NewDate(d.Date)
	for _, h := range d.R.Habits {
		habit := habitToCore(h)
		switch habit.Type {
		case core.HabitTypeWakeUp:
			day.Habits.Sleep = habit
		case core.HabitTypeFitness:
			day.Habits.Fitness = habit
		case core.HabitTypeDeepWork:
			day.Habits.DeepWork = habit
		case core.HabitTypeEatHealthy:
			day.Habits.Eat = habit
		}
		if h.State == string(core.HabitStateDone) {
			day.Habits.Score += 1
		}
		day.Habits.Hs = append(day.Habits.Hs, habit)
	}
	for _, sl := range d.R.SleepLogs {
		day.SleepLogs = append(day.SleepLogs, sleepLogToCore(d, sl))
	}
	for _, fl := range d.R.FitnessLogs {
		day.FitnessLogs = append(day.FitnessLogs, fitnessLogToCore(d, fl))
	}
	for _, wl := range d.R.DeepWorkLogs {
		day.DeepWorkLogs = append(day.DeepWorkLogs, deepWorkLogToCore(d, wl))
	}
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
