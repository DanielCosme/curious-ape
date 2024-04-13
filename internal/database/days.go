package database

import (
	"errors"
	entity2 "github.com/danielcosme/curious-ape/internal/entity"
	"time"
)

type Day interface {
	Create(*entity2.Day) error
	Update(*entity2.Day, ...entity2.DayJoin) (*entity2.Day, error)
	Get(entity2.DayFilter, ...entity2.DayJoin) (*entity2.Day, error)
	Find(entity2.DayFilter, ...entity2.DayJoin) ([]*entity2.Day, error)
}

func ExecuteDaysPipeline(days []*entity2.Day, joins ...entity2.DayJoin) error {
	// TODO implement async pipeline execution
	if !(len(days) > 0) {
		return nil
	}

	for _, j := range joins {
		if err := j(days); err != nil {
			return err
		}
	}
	return nil
}

func DaysPipeline(m *Repository) []entity2.DayJoin {
	return []entity2.DayJoin{
		DaysJoinHabits(m),
		DaysJoinSleepLogs(m),
		DaysJoinFitnessLogs(m),
	}
}

func DaysJoinHabits(m *Repository) entity2.DayJoin {
	return func(days []*entity2.Day) error {
		if len(days) > 0 {
			hs, err := m.Habits.Find(entity2.HabitFilter{DayID: DayToIDs(days)}, HabitsJoinCategories(m))
			if err != nil {
				return err
			}

			habitsByDateID := map[int][]*entity2.Habit{}
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

func DaysJoinSleepLogs(m *Repository) entity2.DayJoin {
	return func(days []*entity2.Day) error {
		if len(days) > 0 {
			sleepLogs, err := m.SleepLogs.Find(entity2.SleepLogFilter{DayID: DayToIDs(days)})
			if err != nil {
				return err
			}

			sleepLogsByDateID := map[int][]*entity2.SleepLog{}
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

func DaysJoinFitnessLogs(m *Repository) entity2.DayJoin {
	return func(days []*entity2.Day) error {
		if len(days) > 0 {
			fitnessLogs, err := m.FitnessLogs.Find(entity2.FitnessLogFilter{DayID: DayToIDs(days)})
			if err != nil {
				return err
			}

			fitnessLogsByDateID := map[int][]*entity2.FitnessLog{}
			for _, log := range fitnessLogs {
				fitnessLogsByDateID[log.DayID] = append(fitnessLogsByDateID[log.DayID], log)
			}

			for _, d := range days {
				d.FitnessLogs = fitnessLogsByDateID[d.ID]
			}
		}
		return nil
	}
}

func DayToIDs(days []*entity2.Day) []int {
	ids := []int{}
	for _, d := range days {
		ids = append(ids, d.ID)
	}
	return ids
}

func DayToMapByISODate(days []*entity2.Day) map[string]*entity2.Day {
	mapDays := map[string]*entity2.Day{}
	for _, d := range days {
		mapDays[entity2.FormatDate(d.Date)] = d
	}
	return mapDays
}

func DayCreate(db *Repository, d *entity2.Day) (*entity2.Day, error) {
	if err := db.Days.Create(d); err != nil {
		return nil, err
	}

	return db.Days.Get(entity2.DayFilter{IDs: []int{d.ID}})
}

func DayGetOrCreate(db *Repository, date time.Time) (*entity2.Day, error) {
	d, err := db.Days.Get(entity2.DayFilter{Dates: []time.Time{date}})
	if err != nil && !errors.Is(err, ErrNotFound) {
		return nil, err
	}
	if d == nil {
		// if it does not exist, create new and return.
		d, err = DayCreate(db, &entity2.Day{Date: date})
		if err != nil {
			return nil, err
		}
	}

	return d, nil
}
