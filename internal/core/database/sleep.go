package database

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
)

type SleepLog interface {
	Create(*entity.SleepLog) error
	Update(*entity.SleepLog, ...entity.SleepLogJoin) (*entity.SleepLog, error)
	Get(entity.SleepLogFilter, ...entity.SleepLogJoin) (*entity.SleepLog, error)
	Find(entity.SleepLogFilter, ...entity.SleepLogJoin) ([]*entity.SleepLog, error)
	Delete(id int) error
}

func ExecuteSleepLogPipeline(ssl []*entity.SleepLog, hjs ...entity.SleepLogJoin) error {
	if !(len(ssl) > 0) {
		return nil
	}

	for _, hj := range hjs {
		if err := hj(ssl); err != nil {
			return err
		}
	}
	return nil
}

func SleepLogsPipeline(m *Repository) []entity.SleepLogJoin {
	return []entity.SleepLogJoin{
		SleepLogsJoinDay(m),
	}
}

func SleepLogsJoinDay(m *Repository) entity.SleepLogJoin {
	return func(sls []*entity.SleepLog) error {
		if len(sls) > 0 {
			days, err := m.Days.Find(entity.DayFilter{IDs: SleepToDayIDs(sls)})
			if err != nil {
				return err
			}

			daysMap := map[int]*entity.Day{}
			for _, d := range days {
				daysMap[d.ID] = d
			}

			for _, h := range sls {
				h.Day = daysMap[h.DayID]
			}
		}
		return nil
	}
}

func SleepToDayIDs(sls []*entity.SleepLog) []int {
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

func SleepToIDs(hs []*entity.SleepLog) []int {
	IDs := []int{}
	mapHabitIDs := map[int]int{}
	for _, h := range hs {
		if _, ok := mapHabitIDs[h.ID]; !ok {
			IDs = append(IDs, h.ID)
			mapHabitIDs[h.ID] = h.ID
		}
	}
	return IDs
}
