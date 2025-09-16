package persistence

import (
	"context"
	"github.com/danielcosme/curious-ape/database/gen/dberrors"
	"github.com/danielcosme/curious-ape/database/gen/models"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
)

type DeepWorkLogs struct {
	db bob.DB
}

func (dw *DeepWorkLogs) Upsert(s *models.DeepWorkLogSetter) (log *models.DeepWorkLog, err error) {
	workLog, err := models.DeepWorkLogs.Insert(s).One(context.Background(), dw.db)
	if err == nil {
		return workLog, nil
	}

	if dberrors.DeepWorkLogErrors.ErrUniqueSqliteAutoindexDeepWorkLog2.Is(err) {
		workLog, err = dw.Get(DeepWorkLogParams{
			DayID:  s.DayID.GetOrZero(),
			Origin: core.OriginLog(s.Origin.GetOrZero()),
		})
		if err != nil {
			return nil, err
		}
		// s.ID = omit.From(workLog.ID)
		err = workLog.Update(context.Background(), dw.db, s)
		if err == nil {
			return workLog, nil
		}
	}
	return nil, catchDBErr("work: upsert", err)
}

func (dw *DeepWorkLogs) Get(p DeepWorkLogParams) (*models.DeepWorkLog, error) {
	workLog, err := p.BuildQuery().One(context.Background(), dw.db)
	if err != nil {
		return nil, catchDBErr("work logs: get", err)
	}
	return workLog, nil
}

type DeepWorkLogParams struct {
	ID     int64
	DayID  int64
	Origin core.OriginLog
}

func (f DeepWorkLogParams) BuildQuery() *sqlite.ViewQuery[*models.DeepWorkLog, models.DeepWorkLogSlice] {
	q := models.DeepWorkLogs.Query()
	if f.ID > 0 {
		q.Apply(models.SelectWhere.DeepWorkLogs.ID.EQ(f.ID))
	}
	if f.DayID > 0 {
		q.Apply(models.SelectWhere.DeepWorkLogs.DayID.EQ(f.DayID))
	}
	if f.Origin != "" {
		q.Apply(models.SelectWhere.DeepWorkLogs.Origin.EQ(string(f.Origin)))
	}
	return q
}
