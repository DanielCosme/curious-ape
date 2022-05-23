package datasource

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/core/repository"
)

func HabitsPipeline(m *repository.Models) []entity.HabitJoin {
	return []entity.HabitJoin{
		HabitsJoinDay(m),
		HabitsJoinCategories(m),
	}
}

func HabitsJoinDay(m *repository.Models) entity.HabitJoin {
	return func(hs []*entity.Habit) error {
		days, err := m.Days.Find(entity.DayFilter{IDs: m.Habits.ToDayIDs(hs)})
		if err != nil {
			return err
		}

		daysMap := map[int]*entity.Day{}
		for _, d := range days {
			daysMap[d.ID] = d
		}

		for _, h := range hs {
			h.Day = daysMap[h.DayID]
		}

		return nil
	}
}

func HabitsJoinCategories(m *repository.Models) entity.HabitJoin {
	return func(hs []*entity.Habit) error {
		cts, err := m.Habits.FindHabitCategories(entity.HabitFilter{CategoryIDs: m.Habits.ToCategoryIDs(hs)})
		if err != nil {
			return err
		}

		ctsMap := map[int]*entity.HabitCategory{}
		for _, c := range cts {
			ctsMap[c.ID] = c
		}

		for _, h := range hs {
			h.Category = ctsMap[h.CategoryID]
		}
		return nil
	}
}
