package persistence

import (
	"context"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"git.danicos.dev/daniel/curious-ape/database/gen/dberrors"
	"git.danicos.dev/daniel/curious-ape/database/gen/models"
	"git.danicos.dev/daniel/curious-ape/pkg/core"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
)

type FitnessLogs struct {
	db bob.DB
}

func (fls *FitnessLogs) Upsert(params core.FitnessLog) (fl core.FitnessLog, err error) {
	day, err := getDay(params.Date, fls.db)
	if err != nil {
		return fl, catchDBErr("fitness logs: upsert: get day", err)
	}
	setter := &models.FitnessLogSetter{
		DayID:     ID(day.ID),
		Title:     omit.From(params.Title),
		StartTime: omit.From(params.StartTime),
		EndTime:   omit.From(params.EndTime),
		Note:      omit.From(params.Note),
		Type:      omit.From(string(params.FitnessType)),
		Origin:    omit.From(string(params.Origin)),
		Raw:       omitnull.From(string(params.Raw)),
	}
	bobFitnessLog, err := models.FitnessLogs.Insert(setter).One(context.Background(), fls.db)
	if err != nil {
		if dberrors.FitnessLogErrors.ErrUniqueSqliteAutoindexFitnessLog1.Is(err) {
			bobFitnessLog, err = fls.Get(FitnessLogParams{
				DayID:     setter.DayID.GetOrZero(),
				StartTime: setter.StartTime.GetOrZero(),
			})
			if err != nil {
				return fl, err
			}
			err = bobFitnessLog.Update(context.Background(), fls.db, setter)
		} else {
			return fl, catchDBErr("fitness: upsert", err)
		}
	}
	return fitnessLogToCore(day, bobFitnessLog), err
}

func (dw *FitnessLogs) Get(p FitnessLogParams) (*models.FitnessLog, error) {
	fitnessLog, err := p.BuildQuery().One(context.Background(), dw.db)
	if err != nil {
		return nil, catchDBErr("fitness logs: get", err)
	}
	return fitnessLog, nil
}

func fitnessLogToCore(day *models.Day, bobFl *models.FitnessLog) (fl core.FitnessLog) {
	fl.ID = uint(bobFl.ID)
	fl.Date = core.NewDate(day.Date)
	fl.Title = bobFl.Title
	fl.StartTime = bobFl.StartTime
	fl.EndTime = bobFl.EndTime
	fl.Note = bobFl.Note
	fl.Type = core.TimelineTypeFitness
	fl.FitnessType = core.FitnessLogType(bobFl.Type)
	fl.Origin = core.LogOrigin(bobFl.Origin)
	return
}

type FitnessLogParams struct {
	ID        int64
	DayID     int64
	Origin    core.LogOrigin
	StartTime time.Time
}

func (f FitnessLogParams) BuildQuery() *sqlite.ViewQuery[*models.FitnessLog, models.FitnessLogSlice] {
	q := models.FitnessLogs.Query()
	if f.ID > 0 {
		q.Apply(models.SelectWhere.FitnessLogs.ID.EQ(f.ID))
	}
	if f.DayID > 0 {
		q.Apply(models.SelectWhere.FitnessLogs.DayID.EQ(f.DayID))
	}
	if f.Origin != "" {
		q.Apply(models.SelectWhere.FitnessLogs.Origin.EQ(string(f.Origin)))
	}
	if !f.StartTime.IsZero() {
		q.Apply(models.SelectWhere.FitnessLogs.StartTime.EQ(f.StartTime))
	}
	return q
}
