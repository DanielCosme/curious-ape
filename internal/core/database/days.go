package database

import (
	"errors"
	"time"

	"github.com/danielcosme/curious-ape/internal/core/entity"
)

type Day interface {
	Create(*entity.Day) error
	Update(*entity.Day, ...entity.DayJoin) (*entity.Day, error)
	Get(entity.DayFilter, ...entity.DayJoin) (*entity.Day, error)
	Find(entity.DayFilter, ...entity.DayJoin) ([]*entity.Day, error)
}

func ExecuteDaysPipeline(days []*entity.Day, joins ...entity.DayJoin) error {
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

func DaysPipeline(m *Repository) []entity.DayJoin {
	return []entity.DayJoin{
		DaysJoinHabits(m),
		DaysJoinSleepLogs(m),
		DaysJoinFitnessLogs(m),
	}
}

func DaysJoinHabits(m *Repository) entity.DayJoin {
	return func(days []*entity.Day) error {
		if len(days) > 0 {
			hs, err := m.Habits.Find(entity.HabitFilter{DayID: DayToIDs(days)}, HabitsJoinCategories(m))
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

func DaysJoinSleepLogs(m *Repository) entity.DayJoin {
	return func(days []*entity.Day) error {
		if len(days) > 0 {
			sleepLogs, err := m.SleepLogs.Find(entity.SleepLogFilter{DayID: DayToIDs(days)})
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

func DaysJoinFitnessLogs(m *Repository) entity.DayJoin {
	return func(days []*entity.Day) error {
		if len(days) > 0 {
			fitnessLogs, err := m.FitnessLogs.Find(entity.FitnessLogFilter{DayID: DayToIDs(days)})
			if err != nil {
				return err
			}

			fitnessLogsByDateID := map[int][]*entity.FitnessLog{}
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

func DayToIDs(days []*entity.Day) []int {
	ids := []int{}
	for _, d := range days {
		ids = append(ids, d.ID)
	}
	return ids
}

func DayToMapByISODate(days []*entity.Day) map[string]*entity.Day {
	mapDays := map[string]*entity.Day{}
	for _, d := range days {
		mapDays[entity.FormatDate(d.Date)] = d
	}
	return mapDays
}

func DayCreate(db *Repository, d *entity.Day) (*entity.Day, error) {
	if err := db.Days.Create(d); err != nil {
		return nil, err
	}

	return db.Days.Get(entity.DayFilter{IDs: []int{d.ID}})
}

func DayGetOrCreate(db *Repository, date time.Time) (*entity.Day, error) {
	d, err := db.Days.Get(entity.DayFilter{Dates: []time.Time{date}})
	if err != nil && !errors.Is(err, ErrNotFound) {
		return nil, err
	}
	if d == nil {
		// if it does not exist, create new and return.
		d, err = DayCreate(db, &entity.Day{Date: date})
		if err != nil {
			return nil, err
		}
	}

	return d, nil
}
