package persistence

import (
	"context"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/danielcosme/curious-ape/database/gen/dberrors"
	"github.com/danielcosme/curious-ape/database/gen/models"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
)

type SleepLogs struct {
	db bob.DB
}

func (fls SleepLogs) Upsert(params core.SleepLogUpsertParams) (sl core.SleepLog, err error) {
	day, err := getDay(params.Date, fls.db)
	if err != nil {
		return sl, catchDBErr("sleep logs: upsert: get day", err)
	}
	s := &models.SleepLogSetter{
		DayID:          omit.From(day.ID),
		StartTime:      omit.From(params.StartTime),
		EndTime:        omit.From(params.EndTime),
		IsMainSleep:    omitnull.From(params.IsMainSleep),
		TotalTimeInBed: omitnull.From(int64(params.TimeInBed)),
		TimeAsleep:     omitnull.From(int64(params.TimeAsleep)),
		Origin:         omit.From(string(params.Origin)),
		Raw:            omitnull.From(params.Raw),
		NOTE:           omitnull.From(params.Note),
	}
	sleepLog, err := models.SleepLogs.Insert(s).One(context.Background(), fls.db)
	if err != nil {
		if dberrors.SleepLogErrors.ErrUniqueSqliteAutoindexSleepLog1.Is(err) {
			ref := s.IsMainSleep.GetOrZero()
			sleepLog, err = fls.Get(SleepLogParams{
				DayID:       s.DayID.GetOrZero(),
				IsMainSleep: &ref,
			})
			if err != nil {
				return
			}

			err = sleepLog.Update(context.Background(), fls.db, s)
			if err != nil {
				return
			}
		} else {
			return
		}
	}

	return sleepLogToCore(day, sleepLog), catchDBErr("sleep: upsert", err)
}

func sleepLogToCore(day *models.Day, s *models.SleepLog) core.SleepLog {
	title := "Nap"
	if s.IsMainSleep.GetOrZero() {
		title = "Main sleep"
	}
	sl := core.SleepLog{
		Date:        core.NewDate(day.Date),
		IsMainSleep: s.IsMainSleep.GetOrZero(),
		TimeAsleep:  time.Duration(s.TimeAsleep.GetOrZero()),
		TimeInBed:   time.Duration(s.TotalTimeInBed.GetOrZero()),
		TimelineLog: core.TimelineLog{
			Title:     title,
			StartTime: s.StartTime,
			EndTime:   s.EndTime,
			Type:      core.TimelineTypeSleep,
			Note:      s.NOTE.GetOrZero(),
		},
	}
	return sl
}

func (sls SleepLogs) Get(p SleepLogParams) (*models.SleepLog, error) {
	sleepLog, err := p.BuildQuery().One(context.Background(), sls.db)
	if err != nil {
		return nil, catchDBErr("sleep logs: get", err)
	}
	return sleepLog, nil
}

type SleepLogParams struct {
	ID          int64
	DayID       int64
	Origin      core.LogOrigin
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
