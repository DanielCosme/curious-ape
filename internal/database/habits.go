package database

import (
	"errors"
	entity2 "github.com/danielcosme/curious-ape/internal/entity"
	"strings"
	"time"
)

type Habit interface {
	// habit
	Create(*entity2.Habit) error
	Update(*entity2.Habit, ...entity2.HabitJoin) (*entity2.Habit, error)
	Get(entity2.HabitFilter, ...entity2.HabitJoin) (*entity2.Habit, error)
	Find(entity2.HabitFilter, ...entity2.HabitJoin) ([]*entity2.Habit, error)
	Delete(id int) error
	// habit log
	CreateHabitLog(*entity2.HabitLog) error
	UpdateHabitLog(*entity2.HabitLog) (*entity2.HabitLog, error)
	GetHabitLog(entity2.HabitLogFilter) (*entity2.HabitLog, error)
	FindHabitLogs(entity2.HabitLogFilter) ([]*entity2.HabitLog, error)
	DeleteHabitLog(id int) error
	// habit category
	GetHabitCategory(entity2.HabitCategoryFilter) (*entity2.HabitCategory, error)
	FindHabitCategories(entity2.HabitCategoryFilter) ([]*entity2.HabitCategory, error)
}

func ExecuteHabitsPipeline(hs []*entity2.Habit, hjs ...entity2.HabitJoin) error {
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

func HabitsPipeline(m *Repository) []entity2.HabitJoin {
	return []entity2.HabitJoin{
		HabitsJoinDay(m),
		HabitsJoinCategories(m),
		HabitsJoinLogs(m),
	}
}

func HabitsJoinDay(m *Repository) entity2.HabitJoin {
	return func(hs []*entity2.Habit) error {
		if len(hs) > 0 {
			days, err := m.Days.Find(entity2.DayFilter{IDs: HabitToDayIDs(hs)})
			if err != nil {
				return err
			}

			daysMap := map[int]*entity2.Day{}
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

func HabitsJoinCategories(m *Repository) entity2.HabitJoin {
	return func(hs []*entity2.Habit) error {
		if len(hs) > 0 {
			cts, err := m.Habits.FindHabitCategories(entity2.HabitCategoryFilter{ID: HabitToCategoryIDs(hs)})
			if err != nil {
				return err
			}

			ctsMap := map[int]*entity2.HabitCategory{}
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

func HabitsJoinLogs(m *Repository) entity2.HabitJoin {
	return func(hs []*entity2.Habit) error {
		if len(hs) > 0 {
			hls, err := m.Habits.FindHabitLogs(entity2.HabitLogFilter{HabitID: HabitToIDs(hs)})
			if err != nil {
				return err
			}

			mapHabits := map[int]*entity2.Habit{}
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

func HabitToDayIDs(hs []*entity2.Habit) []int {
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

func HabitToCategoryIDs(hs []*entity2.Habit) []int {
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

func HabitToIDs(hs []*entity2.Habit) []int {
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

func GetOrCreateHabit(db *Repository, date time.Time, categoryCode string, joins ...entity2.HabitJoin) (*entity2.Habit, error) {
	// Make sure the category exists.
	category, err := db.Habits.GetHabitCategory(entity2.HabitCategoryFilter{Code: []string{strings.ToLower(categoryCode)}})
	if err != nil {
		return nil, err
	}

	day, err := DayGetOrCreate(db, date)
	if err != nil {
		return nil, err
	}

	// First check if the habit already exists.
	h, err := db.Habits.Get(entity2.HabitFilter{DayID: []int{day.ID}, CategoryID: []int{category.ID}}, joins...)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	// If it does not exist we create it.
	if h == nil {
		h = &entity2.Habit{
			DayID:      day.ID,
			CategoryID: category.ID,
			Status:     entity2.HabitStatusNoInfo,
		}
		if err := db.Habits.Create(h); err != nil {
			return nil, err
		}

		h.Day = day
		h.Category = category
	}

	return h, nil
}

func UpsertHabitLog(db *Repository, data *entity2.HabitLog) (string, error) {
	if data.HabitID == 0 {
		return "", errors.New("habit ID cannot be 0")
	}

	hl, err := db.Habits.GetHabitLog(entity2.HabitLogFilter{Origin: []entity2.DataSource{data.Origin}, HabitID: []int{data.HabitID}})
	if err != nil && !errors.Is(err, ErrNotFound) {
		return "", err
	}
	// If it does not exist create it
	if hl == nil {
		return "created", db.Habits.CreateHabitLog(data)
	}

	hl.Origin = data.Origin
	hl.Note = data.Note
	hl.Success = data.Success
	hl.IsAutomated = data.IsAutomated
	_, err = db.Habits.UpdateHabitLog(hl)
	return "update", err
}
