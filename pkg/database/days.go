package database

import (
	"context"
	"github.com/aarondl/opt/omit"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/database/gen/models"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
)

type Days struct {
	db bob.DB
}

func (d *Days) Create(s *models.DaySetter) (*models.Day, error) {
	res, err := models.Days.Insert(s).One(context.Background(), d.db)
	if err != nil {
		return nil, catchDBErr("days: create", err)
	}
	return res, nil
}

func (d *Days) Get(p DayParams) (*models.Day, error) {
	res, err := p.BuildQuery().One(context.Background(), d.db)
	if err != nil {
		return nil, catchDBErr("days: get", err)
	}
	if p.LoadHabits {
		if err := d.LoadHabitRelations(res); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *Days) Find(p DayParams) ([]*models.Day, error) {
	res, err := p.BuildQuery().All(context.Background(), d.db)
	if err != nil {
		return nil, catchDBErr("days: find", err)
	}

	for _, day := range res {
		if p.LoadHabits {
			if err := d.LoadHabitRelations(day); err != nil {
				return nil, err
			}
		}
	}
	return res, nil
}

func (d *Days) GetOrCreate(p DayParams) (*models.Day, error) {
	res, err := d.Get(p)
	if IgnoreIfErrNotFound(err) {
		return res, err
	}
	if res == nil {
		res, err = d.Create(&models.DaySetter{Date: omit.From(p.Date.Time())})
	}
	return res, err
}

func (d *Days) LoadHabitRelations(m *models.Day) error {
	if err := m.R.Habits.LoadHabitHabitCategory(context.Background(), d.db); err != nil {
		return catchDBErr("days: load: habit relations", err)
	}
	return nil
}

type DayParams struct {
	ID         int32
	Date       core.Date
	Dates      core.DateSlice
	R          []Relation
	LoadHabits bool
}

func (f *DayParams) BuildQuery() *sqlite.ViewQuery[*models.Day, models.DaySlice] {
	q := models.Days.Query()
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
			f.LoadHabits = true
		case RelationSleep:
			q.Apply(models.ThenLoadDaySleepLogs())
		case RelationFitness:
			q.Apply(models.ThenLoadDayFitnessLogs())
		case RelationWork:
			q.Apply(models.ThenLoadDayDeepWorkLogs())
		}
	}
	return q
}

func DayRelations() []Relation {
	return []Relation{
		RelationHabit,
		RelationFitness,
		RelationSleep,
		RelationWork,
	}
}
