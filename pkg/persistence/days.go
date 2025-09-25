package persistence

import (
	"context"
	"log/slog"

	"github.com/aarondl/opt/omit"
	"github.com/danielcosme/curious-ape/database/gen/models"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
)

type Days struct {
	db bob.DB
}

func (d Days) Create(date core.Date) (day core.Day, err error) {
	s := &models.DaySetter{Date: omit.From(date.Time())}
	res, err := models.Days.Insert(s).One(context.Background(), d.db)
	return dayToCore(res), err
}

func (d Days) Get(p core.DayParams) (day core.Day, err error) {
	res, err := BuildDayQuery(p).One(context.Background(), d.db)
	if err != nil {
		return day, catchDBErr("days: get", err)
	}
	err = d.LoadHabitRelations(res)
	return dayToCore(res), err
}

func (d Days) GetOrCreate(p core.DayParams) (day core.Day, err error) {
	day, err = d.Get(p)
	if core.IfErrNNotFound(err) {
		return
	}
	if day.IsZero() {
		return d.Create(p.Date)
	}
	return
}

func (d Days) Find(p core.DayParams) (days []core.Day, err error) {
	res, err := BuildDayQuery(p).All(context.Background(), d.db)
	if err != nil {
		return days, catchDBErr("days: find", err)
	}
	//TODO:optimize this.
	for _, day := range res {
		if err = d.LoadHabitRelations(day); err != nil {
			return
		}
		days = append(days, dayToCore(day))
	}
	return
}

func (d *Days) LoadHabitRelations(m *models.Day) error {
	if err := m.R.Habits.LoadDay(context.Background(), d.db); err != nil {
		return catchDBErr("days: load: habit relations", err)
	}
	if err := m.R.Habits.LoadHabitCategory(context.Background(), d.db); err != nil {
		return catchDBErr("days: load: habit relations", err)
	}
	return nil
}

func dayToCore(d *models.Day) (day core.Day) {
	if d == nil {
		slog.Error("dayToCore: day is nil")
		return
	}
	day.ID = uint(d.ID)
	day.Date = core.NewDate(d.Date)
	for _, h := range d.R.Habits {
		day.Habits = append(day.Habits, habitToCore(h))
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
	if !f.Lite {
		q.Apply(
			models.SelectThenLoad.Day.Habits(),
			models.SelectThenLoad.Day.SleepLogs(),
			models.SelectThenLoad.Day.FitnessLogs(),
			models.SelectThenLoad.Day.DeepWorkLogs(),
		)
	}
	return q
}
