package database

import "github.com/danielcosme/curious-ape/internal/core/entity"

type FitnessLog interface {
	Create(*entity.FitnessLog) error
	Update(*entity.FitnessLog, ...entity.FitnessLogJoin) (*entity.FitnessLog, error)
	Get(entity.FitnessLogFilter, ...entity.FitnessLogJoin) (*entity.FitnessLog, error)
	Find(entity.FitnessLogFilter, ...entity.FitnessLogJoin) ([]*entity.FitnessLog, error)
	Delete(id int) error
}

func ExecuteFitnessLogPipeline(fls []*entity.FitnessLog, fjs ...entity.FitnessLogJoin) error {
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

func FitnessLogsJoinDay(m *Repository) entity.FitnessLogJoin {
	return func(fls []*entity.FitnessLog) error {
		if len(fls) > 0 {
			days, err := m.Days.Find(entity.DayFilter{IDs: FitnessToDayIDs(fls)})
			if err != nil {
				return err
			}

			daysMap := map[int]*entity.Day{}
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

func FitnessToDayIDs(sls []*entity.FitnessLog) []int {
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
