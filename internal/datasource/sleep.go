package datasource

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/core/repository"
)

func SleepLogsPipeline(m *repository.Models) []entity.SleepLogJoin {
	return []entity.SleepLogJoin{
		SleepLogsJoinDay(m),
	}
}

func SleepLogsJoinDay(m *repository.Models) entity.SleepLogJoin {
	return func(sls []*entity.SleepLog) error {
		if len(sls) > 0 {
			days, err := m.Days.Find(entity.DayFilter{IDs: m.SleepLogs.ToDayIDs(sls)})
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
