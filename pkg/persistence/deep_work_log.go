package persistence

import (
	"context"
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/gen/bob/dberrors"
	"danicos.dev/daniel/curious-ape/pkg/gen/bob/models"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
)

type DeepWorkLogs struct {
	db bob.DB
}

func (dw *DeepWorkLogs) Upsert(params core.DeepWorkLog) (log core.DeepWorkLog, err error) {
	day, err := GetDay(params.Date, dw.db)
	if err != nil {
		return log, CatchDBErr("fitness logs: upsert: get day", err)
	}
	setter := &models.DeepWorkLogSetter{
		DayID:     ID(day.ID),
		Title:     omit.From(params.Title),
		StartTime: omit.From(params.StartTime),
		EndTime:   omit.From(params.EndTime),
		Note:      omit.From(params.Note),
		Origin:    omit.From(string(params.Origin)),
		Raw:       omitnull.From(string(params.Raw)),
	}
	bobLog, err := models.DeepWorkLogs.Insert(setter).One(context.Background(), dw.db)
	if err != nil {
		if dberrors.DeepWorkLogErrors.ErrUniqueSqliteAutoindexDeepWorkLog1.Is(err) {
			bobLog, err = dw.Get(DeepWorkLogParams{
				DayID:     day.ID,
				StartTime: params.StartTime,
			})
			if err != nil {
				return
			}
			err = bobLog.Update(context.Background(), dw.db, setter)
		} else {
			return
		}
	}

	return deepWorkLogToCore(day, bobLog), CatchDBErr("work: upsert", err)
}

func deepWorkLogToCore(day *models.Day, bob *models.DeepWorkLog) (log core.DeepWorkLog) {
	log.ID = uint(bob.ID)
	log.Date = core.NewDate(day.Date)
	log.Title = bob.Title
	log.StartTime = bob.StartTime
	log.EndTime = bob.EndTime
	log.Note = bob.Note
	log.Type = core.TimelineTypeDeepWork
	return
}

func (dw *DeepWorkLogs) Get(p DeepWorkLogParams) (*models.DeepWorkLog, error) {
	workLog, err := p.BuildQuery().One(context.Background(), dw.db)
	if err != nil {
		return nil, CatchDBErr("work logs: get", err)
	}
	return workLog, nil
}

type DeepWorkLogParams struct {
	ID        int64
	DayID     int64
	StartTime time.Time
}

func (f DeepWorkLogParams) BuildQuery() *sqlite.ViewQuery[*models.DeepWorkLog, models.DeepWorkLogSlice] {
	q := models.DeepWorkLogs.Query()
	if f.ID > 0 {
		q.Apply(models.SelectWhere.DeepWorkLogs.ID.EQ(f.ID))
	}
	if f.DayID > 0 {
		q.Apply(models.SelectWhere.DeepWorkLogs.DayID.EQ(f.DayID))
	}
	if !f.StartTime.IsZero() {
		q.Apply(models.SelectWhere.DeepWorkLogs.StartTime.EQ(f.StartTime))
	}
	return q
}
