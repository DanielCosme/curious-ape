package database

import (
	"context"
	"errors"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database/gen/models"
	"github.com/stephenafamo/bob"
)

type FitnessLogs struct {
	db bob.DB
}

func (fls *FitnessLogs) Upsert(fitnessLog core.FitnessLog) (core.FitnessLog, error) {
	if fitnessLog.DayID == 0 {
		return core.FitnessLog{}, errors.New("day ID cannot be 0")
	}
	fl, err := models.FitnessLogs.Upsert(
		context.Background(),
		fls.db,
		true,
		[]string{"day_id", "start_time"},
		nil,
		fromFitnessLogToSetter(fitnessLog),
	)
	return fitnessLogToCore(fl), err
}

func fromFitnessLogToSetter(fl core.FitnessLog) *models.FitnessLogSetter {
	return &models.FitnessLogSetter{
		ID:        omit.FromCond(fl.ID, fl.ID > 0),
		DayID:     omit.From(fl.DayID),
		Date:      omit.From(fl.Date.Time()),
		StartTime: omit.From(fl.StartTime),
		EndTime:   omit.From(fl.EndTime),
		Type:      omit.From(fl.Type),
		Title:     omit.From(fl.Title),
		Origin:    omit.From(string(fl.Origin)),
		Note:      omitnull.From(fl.Note),
		Raw:       omitnull.From(fl.Raw),
	}
}

func fitnessLogToCore(m *models.FitnessLog) (fl core.FitnessLog) {
	if m == nil {
		return fl
	}
	fl.ID = m.ID
	fl.DayID = m.DayID
	fl.Date = core.NewDate(m.Date)
	fl.StartTime = m.StartTime
	fl.EndTime = m.EndTime
	fl.Type = m.Type
	fl.Title = m.Title
	fl.Origin = core.Integration(m.Origin)
	fl.Note = m.Note.GetOrZero()
	if !m.Raw.IsNull() {
		fl.Raw = m.Raw.MustGet()
	}
	return
}
