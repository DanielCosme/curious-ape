package database

import (
	entity2 "github.com/danielcosme/curious-ape/internal/entity"
)

type FitnessLog interface {
	Create(*entity2.FitnessLog) error
	Update(*entity2.FitnessLog, ...entity2.FitnessLogJoin) (*entity2.FitnessLog, error)
	Get(entity2.FitnessLogFilter, ...entity2.FitnessLogJoin) (*entity2.FitnessLog, error)
	Find(entity2.FitnessLogFilter, ...entity2.FitnessLogJoin) ([]*entity2.FitnessLog, error)
	Delete(id int) error
}

func ExecuteFitnessLogPipeline(fls []*entity2.FitnessLog, fjs ...entity2.FitnessLogJoin) error {
	if !(len(fls) > 0) {
		return nil
	}

	for _, fj := range fjs {
		if err := fj(fls); err != nil {
			return err
		}
	}
	return nil
}

func FitnessLogsJoinDay(m *Repository) entity2.FitnessLogJoin {
	return func(fls []*entity2.FitnessLog) error {
		if len(fls) > 0 {
			days, err := m.Days.Find(entity2.DayFilter{IDs: FitnessToDayIDs(fls)})
			if err != nil {
				return err
			}

			daysMap := map[int]*entity2.Day{}
			for _, d := range days {
				daysMap[d.ID] = d
			}

			for _, h := range fls {
				h.Day = daysMap[h.DayID]
			}
		}
		return nil
	}
}

func FitnessToDayIDs(sls []*entity2.FitnessLog) []int {
	dayIDs := []int{}
	dayIDsMap := map[int]int{}
	for _, h := range sls {
		if _, ok := dayIDsMap[h.DayID]; !ok {
			dayIDs = append(dayIDs, h.DayID)
			dayIDsMap[h.DayID] = h.DayID
		}
	}
	return dayIDs
}
