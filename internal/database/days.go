package database

import (
	"context"
	"errors"
	"github.com/aarondl/opt/omit"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database/gen/models"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"log/slog"
)

type Days struct {
	db bob.DB
}

func (d *Days) Query() *sqlite.ViewQuery[*models.Day, models.DaySlice] {
	return models.Days.Query(context.Background(), d.db)
}

func (d *Days) GetOrCreate(f DayF) (*core.Day, error) {
	res, err := d.Get(f)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return nil, err
	}
	if res == nil {
		res, err = d.Create(models.DaySetter{Date: omit.From(f.Date.Time())})
	}
	return res, err
}

func (d *Days) Create(s models.DaySetter) (*core.Day, error) {
	res, err := models.Days.Insert(context.Background(), d.db, &s)
	if err != nil {
		return nil, catchErr("create day", err)
	}
	return dayToCore(res), nil
}

func (d *Days) Find(f DayF) ([]*core.Day, error) {
	q := f.BuildQuery(d.db)
	res, err := q.All()
	if err != nil {
		return nil, catchErr("find days", err)
	}
	bb, _, _ := q.Build()
	slog.Debug(bb)
	if f.WithAll {
		err := res.LoadDayHabits(context.Background(), d.db)
		if err != nil {
			return nil, err
		}
	}
	return daysToCore(res), nil
}

func (d *Days) Get(f DayF) (*core.Day, error) {
	q := f.BuildQuery(d.db)
	res, err := q.One()
	if err != nil {
		return nil, catchErr("get day", err)
	}
	if f.WithAll {
		err := res.LoadDayHabits(context.Background(), d.db)
		if err != nil {
			return nil, err
		}
	}
	return dayToCore(res), nil
}

func dayToCore(m *models.Day) *core.Day {
	d := &core.Day{
		ID:   m.ID,
		Date: core.NewDate(m.Date),
	}
	// TODO: populate habits.
	// for _, h := range m.R.Habits {
	//
	// }
	return d
}

func daysToCore(ds models.DaySlice) []*core.Day {
	res := make([]*core.Day, len(ds))
	for idx, d := range ds {
		res[idx] = dayToCore(d)
	}
	return res
}
