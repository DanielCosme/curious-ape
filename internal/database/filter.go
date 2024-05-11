package database

import (
	"context"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database/gen/models"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
)

type DayF struct {
	Date    core.Date
	Dates   core.DateSlice
	WithAll bool
}

func (f DayF) BuildQuery(exec bob.Executor) *sqlite.ViewQuery[*models.Day, models.DaySlice] {
	q := models.Days.Query(context.Background(), exec)

	if !f.Date.Time().IsZero() {
		q.Apply(models.SelectWhere.Days.Date.EQ(f.Date.Time()))
	}

	if len(f.Dates) > 0 {
		q.Apply(models.SelectWhere.Days.Date.In(f.Dates.ToTimeSlice()...))
	}

	return q
}
