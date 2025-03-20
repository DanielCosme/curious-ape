package database

import (
	"context"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/database/gen/models"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"time"
)

type SleepLogs struct {
	db bob.DB
}

func (fls *SleepLogs) Upsert(s *models.SleepLogSetter) (*models.SleepLog, error) {
	sleepLog, err := models.SleepLogs.Insert(s).One(context.Background(), fls.db)
	if err == nil {
		return sleepLog, nil
	}

	if models.SleepLogErrors.ErrUniqueDayIdAndIsMainSleep.Is(err) {
		ref := s.IsMainSleep.GetOrZero()
		sleepLog, err = fls.Get(SleepLogParams{
			DayID:       s.DayID.GetOrZero(),
			IsMainSleep: &ref,
		})
		if err != nil {
			return nil, err
		}

		err = sleepLog.Update(context.Background(), fls.db, s)
		if err == nil {
			return sleepLog, nil
		}
	}
	return nil, catchDBErr("sleep: upsert", err)
}

func (sls *SleepLogs) Get(p SleepLogParams) (*models.SleepLog, error) {
	sleepLog, err := p.BuildQuery().One(context.Background(), sls.db)
	if err != nil {
		return nil, catchDBErr("sleep logs: get", err)
	}
	return sleepLog, nil
}

type SleepLogParams struct {
	ID          int32
	DayID       int32
	Origin      core.OriginLog
	IsMainSleep *bool
}

func (f SleepLogParams) BuildQuery() *sqlite.ViewQuery[*models.SleepLog, models.SleepLogSlice] {
	q := models.SleepLogs.Query()
	if f.ID > 0 {
		q.Apply(models.SelectWhere.SleepLogs.ID.EQ(f.ID))
	}
	if f.DayID > 0 {
		q.Apply(models.SelectWhere.SleepLogs.DayID.EQ(f.DayID))
	}
	if f.Origin != "" {
		q.Apply(models.SelectWhere.SleepLogs.Origin.EQ(string(f.Origin)))
	}
	if f.IsMainSleep != nil {
		q.Apply(models.SelectWhere.SleepLogs.IsMainSleep.EQ(*f.IsMainSleep))
	}
	return q
}

func toDuration(i int) time.Duration {
	return time.Duration(i) * time.Minute
}
