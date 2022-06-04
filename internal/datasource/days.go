package datasource

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/core/repository"
)

func DaysPipeline(m *repository.Models) []entity.DayJoin {
	return []entity.DayJoin{
		DaysJoinHabits(m),
		DaysJoinSleepLogs(m),
	}
}

func DaysJoinHabits(m *repository.Models) entity.DayJoin {
	return func(days []*entity.Day) error {
		if len(days) > 0 {
			hs, err := m.Habits.Find(entity.HabitFilter{DayID: m.Days.ToIDs(days)}, HabitsJoinCategories(m))
			if err != nil {
				return err
			}

			habitsByDateID := map[int][]*entity.Habit{}
			for _, h := range hs {
				habitsByDateID[h.DayID] = append(habitsByDateID[h.DayID], h)
			}

			for _, d := range days {
				d.Habits = habitsByDateID[d.ID]
			}
		}
		return nil
	}
}

func DaysJoinSleepLogs(m *repository.Models) entity.DayJoin {
	return func(days []*entity.Day) error {
		if len(days) > 0 {
			sleepLogs, err := m.SleepLogs.Find(entity.SleepLogFilter{DayID: m.Days.ToIDs(days)})
			if err != nil {
				return err
			}

			sleepLogsByDateID := map[int][]*entity.SleepLog{}
			for _, log := range sleepLogs {
				sleepLogsByDateID[log.DayID] = append(sleepLogsByDateID[log.DayID], log)
			}

			for _, d := range days {
				d.SleepLogs = sleepLogsByDateID[d.ID]
			}
		}
		return nil
	}
}
