package datasource

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/core/repository"
)

func DaysPipeline(m *repository.Models) []entity.DayJoin {
	return []entity.DayJoin{
		DaysJoinHabits(m),
	}
}

func DaysJoinHabits(m *repository.Models) entity.DayJoin {
	return func(days []*entity.Day) error {
		hs, err := m.Habits.Find(entity.HabitFilter{DayIDs: m.Days.ToIDs(days)}, HabitsJoinCategories(m))
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

		return nil
	}
}
