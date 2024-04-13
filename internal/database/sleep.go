package database

import (
	entity2 "github.com/danielcosme/curious-ape/internal/entity"
)

type SleepLog interface {
	Create(*entity2.SleepLog) error
	Update(*entity2.SleepLog, ...entity2.SleepLogJoin) (*entity2.SleepLog, error)
	Get(entity2.SleepLogFilter, ...entity2.SleepLogJoin) (*entity2.SleepLog, error)
	Find(entity2.SleepLogFilter, ...entity2.SleepLogJoin) ([]*entity2.SleepLog, error)
	Delete(id int) error
}

func ExecuteSleepLogPipeline(ssl []*entity2.SleepLog, hjs ...entity2.SleepLogJoin) error {
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

func SleepLogsPipeline(m *Repository) []entity2.SleepLogJoin {
	return []entity2.SleepLogJoin{
		SleepLogsJoinDay(m),
	}
}

func SleepLogsJoinDay(m *Repository) entity2.SleepLogJoin {
	return func(sls []*entity2.SleepLog) error {
		if len(sls) > 0 {
			days, err := m.Days.Find(entity2.DayFilter{IDs: SleepToDayIDs(sls)})
			if err != nil {
				return err
			}

			daysMap := map[int]*entity2.Day{}
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

func SleepToDayIDs(sls []*entity2.SleepLog) []int {
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

func SleepToIDs(hs []*entity2.SleepLog) []int {
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
