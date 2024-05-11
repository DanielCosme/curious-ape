package database

import (
	"context"
	"github.com/aarondl/opt/omit"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database/gen/models"
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
	return dayToCore(res), nil
}

func (d *Days) Find(p DayParams) ([]core.Day, error) {
	res, err := p.BuildQuery(d.db).All()
	if err != nil {
		return nil, catchErr("find days", err)
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

func dayToCore(m *models.Day) core.Day {
	d := core.Day{
		ID:   m.ID,
		Date: core.NewDate(m.Date),
	}
	d.Habits = habitsToCore(m.R.Habits)
	return d
}

func daysToCore(ds models.DaySlice) []core.Day {
	res := make([]core.Day, len(ds))
	for idx, d := range ds {
		res[idx] = dayToCore(d)
	}
	return res
}
