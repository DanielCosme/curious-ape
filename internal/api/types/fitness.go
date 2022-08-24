package types

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
)

type FitnessLogTransport struct {
	ID        int                   `json:"id"`
	Title     string                `json:"title"`
	Type      entity.FitnessLogType `json:"type"`
	Origin    entity.DataSource     `json:"origin"`
	Date      string                `json:"date"`
	StartTime string                `json:"start_time"`
	EndTime   string                `json:"end_time"`
	Note      string                `json:"note"`
}

func FromFitnessLogToTransport(fl *entity.FitnessLog) *FitnessLogTransport {
	var flt *FitnessLogTransport
	if fl != nil {
		flt = &FitnessLogTransport{
			ID:        fl.ID,
			Type:      fl.Type,
			Title:     fl.Title,
			Origin:    fl.Origin,
			StartTime: fl.StartTime.Format(entity.Timestamp),
			EndTime:   fl.EndTime.Format(entity.Timestamp),
			Note:      fl.Note,
		}

		if fl.Day != nil {
			flt.Date = fl.Day.FormatDate()
		}
	}
	return flt
}

func FromFitnessLogToTransportSlice(fls []*entity.FitnessLog) []*FitnessLogTransport {
	flst := []*FitnessLogTransport{}
	for _, fl := range fls {
		flst = append(flst, FromFitnessLogToTransport(fl))
	}
	return flst
}

func (flt *FitnessLogTransport) ToFitnessLog(day *entity.Day) (*entity.FitnessLog, error) {
	startTime, err := entity.ParseTime(flt.StartTime)
	if err != nil {
		return nil, err
	}
	endTime, err := entity.ParseTime(flt.EndTime)
	if err != nil {
		return nil, err
	}

	return &entity.FitnessLog{
		DayID:     day.ID,
		Type:      flt.Type,
		Title:     flt.Title,
		Date:      day.Date,
		StartTime: startTime,
		EndTime:   endTime,
		Origin:    flt.Origin,
		Note:      flt.Note,
		Day:       day,
	}, nil
}
