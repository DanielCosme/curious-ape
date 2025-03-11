package database

import (
	"context"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/database/gen/models"
	"github.com/stephenafamo/bob"
	"time"
)

type SleepLogs struct {
	db bob.DB
}

func (sls *SleepLogs) Upsert(sleepLog core.SleepLog) (core.SleepLog, error) {
	setter := fromSleepLogCoreToSetter(sleepLog)
	sl, err := models.SleepLogs.Upsert(
		context.Background(),
		sls.db,
		true,
		[]string{"day_id", "is_main_sleep"},
		nil,
		setter,
	)
	if err != nil {
		return core.SleepLog{}, err
	}
	return sleepLogToCore(sl), nil
}

func fromSleepLogCoreToSetter(sl core.SleepLog) *models.SleepLogSetter {
	return &models.SleepLogSetter{
		ID:             omit.FromCond(sl.ID, sl.ID > 0),
		DayID:          omit.From(sl.DayID),
		Date:           omit.From(sl.Date.Time()),
		StartTime:      omit.From(sl.StartTime),
		EndTime:        omit.From(sl.EndTime),
		IsMainSleep:    omitnull.From(sl.IsMainSleep),
		IsAutomated:    omitnull.From(sl.IsAutomated),
		Origin:         omit.From(string(sl.Origin)),
		TotalTimeInBed: omitnull.From(int32(sl.MinutesInBed.Minutes())),
		MinutesAsleep:  omitnull.From(int32(sl.MinutesAsleep.Minutes())),
		MinutesRem:     omitnull.From(int32(sl.MinutesRem.Minutes())),
		MinutesDeep:    omitnull.From(int32(sl.MinutesDeep.Minutes())),
		MinutesLight:   omitnull.From(int32(sl.MinutesLight.Minutes())),
		MinutesAwake:   omitnull.From(int32(sl.MinutesAwake.Minutes())),
		Raw:            omitnull.From(string(sl.Raw)),
	}
}

func sleepLogToCore(m *models.SleepLog) (sl core.SleepLog) {
	sl.ID = m.ID
	sl.DayID = m.DayID
	sl.Date = core.NewDate(m.Date)
	sl.StartTime = m.StartTime
	sl.EndTime = m.EndTime
	sl.IsMainSleep = m.IsMainSleep.GetOrZero()
	sl.IsAutomated = m.IsAutomated.GetOrZero()
	sl.MinutesAsleep = toDuration(int(m.MinutesAsleep.GetOrZero()))
	sl.MinutesAwake = toDuration(int(m.MinutesAwake.GetOrZero()))
	sl.MinutesDeep = toDuration(int(m.MinutesDeep.GetOrZero()))
	sl.MinutesRem = toDuration(int(m.MinutesRem.GetOrZero()))
	sl.MinutesLight = toDuration(int(m.MinutesLight.GetOrZero()))
	sl.MinutesInBed = toDuration(int(m.TotalTimeInBed.GetOrZero()))
	return sl
}

func toDuration(i int) time.Duration {
	return time.Duration(i) * time.Minute
}
