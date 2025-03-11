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

type DeepWorkLogs struct {
	db bob.DB
}

func (dw *DeepWorkLogs) Upsert(log core.DeepWorkLog) (res core.DeepWorkLog, err error) {
	setter := fromDeepWorkCoreToSetter(log)
	dwLog, err := models.DeepWorkLogs.Upsert(
		context.Background(),
		dw.db,
		true,
		[]string{"origin", "day_id"},
		nil,
		setter,
	)
	if err != nil {
		return
	}
	return deepWorkLogToCore(dwLog), nil
}

func fromDeepWorkCoreToSetter(log core.DeepWorkLog) *models.DeepWorkLogSetter {
	return &models.DeepWorkLogSetter{
		ID:          omit.FromCond(log.ID, log.ID > 0),
		DayID:       omit.From(log.DayID),
		Date:        omit.From(log.Date.Time()),
		Seconds:     omit.From(int32(log.Duration.Seconds())),
		IsAutomated: omitnull.From(log.IsAutomated),
		Origin:      omit.From(string(log.Origin)),
	}
}

func deepWorkLogToCore(m *models.DeepWorkLog) (log core.DeepWorkLog) {
	log.ID = m.ID
	log.DayID = m.DayID
	log.Date = core.NewDate(m.Date)
	log.Duration = time.Duration(m.Seconds) * time.Second
	log.Origin = core.Integration(m.Origin)
	log.IsAutomated = m.IsAutomated.GetOrZero()
	return
}
