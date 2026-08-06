package persistence

import (
	"context"
	"log/slog"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/gen/bob/models"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
)

func GetDay(exec bob.Executor, date core.Date) (*models.Day, error) {
	day, err := BuildDayQuery(core.DayParams{Date: date}).One(context.Background(), exec)
	if err != nil {
		slog.Error("Get day: failed", "err", err)
		return nil, err
	}
	return day, nil
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
	return q
}
