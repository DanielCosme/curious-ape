package database

import "github.com/danielcosme/curious-ape/internal/core/entity"

type Habit interface {
	// habit
	Create(*entity.Habit) error
	Update(*entity.Habit, ...entity.HabitJoin) (*entity.Habit, error)
	Get(entity.HabitFilter, ...entity.HabitJoin) (*entity.Habit, error)
	Find(entity.HabitFilter, ...entity.HabitJoin) ([]*entity.Habit, error)
	Delete(id int) error
	// habit log
	CreateHabitLog(*entity.HabitLog) error
	UpdateHabitLog(*entity.HabitLog) (*entity.HabitLog, error)
	GetHabitLog(entity.HabitLogFilter) (*entity.HabitLog, error)
	FindHabitLogs(entity.HabitLogFilter) ([]*entity.HabitLog, error)
	DeleteHabitLog(id int) error
	// habit category
	GetHabitCategory(entity.HabitCategoryFilter) (*entity.HabitCategory, error)
	FindHabitCategories(entity.HabitCategoryFilter) ([]*entity.HabitCategory, error)
}

func ExecuteHabitsPipeline(hs []*entity.Habit, hjs ...entity.HabitJoin) error {
	if !(len(hs) > 0) {
		return nil
	}

	for _, hj := range hjs {
		if err := hj(hs); err != nil {
			return err
		}
	}
	return nil
}

func HabitsPipeline(m *Repository) []entity.HabitJoin {
	return []entity.HabitJoin{
		HabitsJoinDay(m),
		HabitsJoinCategories(m),
		HabitsJoinLogs(m),
	}
}

func HabitsJoinDay(m *Repository) entity.HabitJoin {
	return func(hs []*entity.Habit) error {
		if len(hs) > 0 {
			days, err := m.Days.Find(entity.DayFilter{IDs: HabitToDayIDs(hs)})
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
		}
		return nil
	}
}

func HabitsJoinCategories(m *Repository) entity.HabitJoin {
	return func(hs []*entity.Habit) error {
		if len(hs) > 0 {
			cts, err := m.Habits.FindHabitCategories(entity.HabitCategoryFilter{ID: HabitToCategoryIDs(hs)})
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
		}
		return nil
	}
}

func HabitsJoinLogs(m *Repository) entity.HabitJoin {
	return func(hs []*entity.Habit) error {
		if len(hs) > 0 {
			hls, err := m.Habits.FindHabitLogs(entity.HabitLogFilter{HabitID: HabitToIDs(hs)})
			if err != nil {
				return err
			}

			mapHabits := map[int]*entity.Habit{}
			for _, h := range hs {
				mapHabits[h.ID] = h
			}

			for _, hl := range hls {
				mapHabits[hl.HabitID].Logs = append(mapHabits[hl.HabitID].Logs, hl)
			}
		}
		return nil
	}
}

func HabitToDayIDs(hs []*entity.Habit) []int {
	dayIDs := []int{}
	dayIDsMap := map[int]int{}
	for _, h := range hs {
		if _, ok := dayIDsMap[h.DayID]; !ok {
			dayIDs = append(dayIDs, h.DayID)
			dayIDsMap[h.DayID] = h.DayID
		}
	}
	return dayIDs
}

func HabitToCategoryIDs(hs []*entity.Habit) []int {
	categoryIDs := []int{}
	categoryIDsMap := map[int]int{}
	for _, h := range hs {
		if _, ok := categoryIDsMap[h.CategoryID]; !ok {
			categoryIDs = append(categoryIDs, h.CategoryID)
			categoryIDsMap[h.CategoryID] = h.CategoryID
		}
	}
	return categoryIDs
}

func HabitToIDs(hs []*entity.Habit) []int {
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
