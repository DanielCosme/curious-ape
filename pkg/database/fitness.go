package database

import (
	"context"
	"github.com/aarondl/opt/omit"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/database/gen/models"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"time"
)

type FitnessLogs struct {
	db bob.DB
}

func (fls *FitnessLogs) Upsert(s *models.FitnessLogSetter) (*models.FitnessLog, error) {
	fitnessLog, err := models.FitnessLogs.Insert(s).One(context.Background(), fls.db)
	if err == nil {
		return fitnessLog, nil
	}

	if models.FitnessLogErrors.ErrUniqueDayIdAndStartTime.Is(err) {
		fitnessLog, err = fls.Get(FitnessLogParams{
			DayID:     s.DayID.GetOrZero(),
			StartTime: s.StartTime.GetOrZero(),
		})
		if err != nil {
			return nil, err
		}
		s.ID = omit.From(fitnessLog.ID)
		fitnessLog, err = models.FitnessLogs.Update(s.UpdateMod()).One(context.Background(), fls.db)
		if err == nil {
			return fitnessLog, nil
		}
	}
	return nil, catchDBErr("fitness: upsert", err)
}

func (dw *FitnessLogs) Get(p FitnessLogParams) (*models.FitnessLog, error) {
	fitnessLog, err := p.BuildQuery().One(context.Background(), dw.db)
	if err != nil {
		return nil, catchDBErr("fitness logs: get", err)
	}
	return fitnessLog, nil
}

type FitnessLogParams struct {
	ID        int32
	DayID     int32
	Origin    core.OriginLog
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
