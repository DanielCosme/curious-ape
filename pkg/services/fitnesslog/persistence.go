package fitnesslog

import (
	"context"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/gen/bob/dberrors"
	"danicos.dev/daniel/curious-ape/pkg/gen/bob/models"
	"danicos.dev/daniel/curious-ape/pkg/persistence"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
)

func upsert(db bob.DB, params core.FitnessLog) (fl core.FitnessLog, err error) {
	day, err := persistence.GetDay(db, params.Date)
	if err != nil {
		return fl, persistence.CatchDBErr("fitness logs: upsert: get day", err)
	}
	setter := &models.FitnessLogSetter{
		DayID:     persistence.ID(day.ID),
		Title:     omit.From(params.Title),
		StartTime: omit.From(params.StartTime),
		EndTime:   omit.From(params.EndTime),
		Note:      omit.From(params.Note),
		Type:      omit.From(string(params.FitnessType)),
		Origin:    omit.From(string(params.Origin)),
		Raw:       omitnull.From(string(params.Raw)),
	}
	bobFitnessLog, err := models.FitnessLogs.Insert(setter).One(context.Background(), db)
	if err != nil {
		if dberrors.FitnessLogErrors.ErrUniqueSqliteAutoindexFitnessLog1.Is(err) {
			bobFitnessLog, err = get(db, core.FitnessLogParams{
				DayID:     setter.DayID.GetOrZero(),
				StartTime: setter.StartTime.GetOrZero(),
			})
			if err != nil {
				return fl, err
			}
			err = bobFitnessLog.Update(context.Background(), db, setter)
		} else {
			return fl, persistence.CatchDBErr("fitness: upsert", err)
		}
	}
	return fitnessLogToCore(day, bobFitnessLog), err
}

func get(db bob.DB, params core.FitnessLogParams) (*models.FitnessLog, error) {
	fitnessLog, err := BuildQuery(params).One(context.Background(), db)
	if err != nil {
		return nil, persistence.CatchDBErr("fitness logs: get", err)
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

func BuildQuery(f core.FitnessLogParams) *sqlite.ViewQuery[*models.FitnessLog, models.FitnessLogSlice] {
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
