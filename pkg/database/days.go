package database

import (
	"context"
	"github.com/aarondl/opt/omit"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/database/gen/models"
	"github.com/stephenafamo/bob"
)

type Days struct {
	db bob.DB
}

func DayRelations() []Relation {
	return []Relation{RelationHabit, RelationFitness, RelationSleep}
}

func (d *Days) Create(s models.DaySetter) (core.Day, error) {
	res, err := models.Days.Insert(context.Background(), d.db, &s)
	if err != nil {
		return core.Day{}, catchErr("create day", err)
	}
	return dayToCore(res), nil
}

func (d *Days) Get(p DayParams) (core.Day, error) {
	res, err := p.BuildQuery(d.db).One()
	if err != nil {
		return core.Day{}, catchErr("get day", err)
	}
	if err := d.LoadHabitRelations(res); err != nil {
		return core.Day{}, err
	}
	return dayToCore(res), nil
}

func (d *Days) Find(p DayParams) ([]core.Day, error) {
	res, err := p.BuildQuery(d.db).All()
	if err != nil {
		return nil, catchErr("find days", err)
	}
	for _, day := range res {
		if err := d.LoadHabitRelations(day); err != nil {
			return nil, err
		}
	}
	return daysToCore(res), nil
}

func (d *Days) GetOrCreate(p DayParams) (core.Day, error) {
	res, err := d.Get(p)
	if IfNotFoundErr(err) {
		return res, err
	}
	if res.IsZero() {
		res, err = d.Create(models.DaySetter{Date: omit.From(p.Date.Time())})
	}
	return res, err
}

func (d *Days) LoadHabitRelations(m *models.Day) error {
	ctx := context.Background()
	if err := m.R.Habits.LoadHabitHabitCategory(ctx, d.db); err != nil {
		return err
	}
	if err := m.R.Habits.LoadHabitHabitLogs(ctx, d.db); err != nil {
		return err
	}
	return nil
}

func dayToCore(m *models.Day) core.Day {
	day := core.Day{
		ID:   m.ID,
		Date: core.NewDate(m.Date),
	}
	day.Habits = habitsToCore(m.R.Habits)
	return day
}

func daysToCore(ds models.DaySlice) []core.Day {
	res := make([]core.Day, len(ds))
	for idx, day := range ds {
		res[idx] = dayToCore(day)
	}
	return res
}
