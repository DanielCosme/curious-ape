package application

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/core/database"
	"github.com/danielcosme/curious-ape/sdk/errors"
	"time"
)

func (a *App) HabitCreate(day *entity.Day, data *entity.Habit) (*entity.Habit, error) {
	habitCategory, err := a.db.Habits.GetHabitCategory(entity.HabitCategoryFilter{ID: []int{data.CategoryID}})
	if err != nil {
		return nil, err
	}

	habit, err := a.getOrCreateHabit(day.ID, habitCategory.ID)
	if err != nil {
		return nil, err
	}

	// TODO hardcode origins for habit logs and create a validator for it
	// Create the habit log
	for _, dataLog := range data.Logs {
		hl, err := a.db.Habits.GetHabitLog(entity.HabitLogFilter{Origin: []entity.HabitOrigin{dataLog.Origin}, HabitID: []int{habit.ID}})
		if err != nil && !errors.Is(err, database.ErrNotFound) {
			return nil, err
		}

		// if it does not exist create it
		if hl == nil {
			dataLog.HabitID = habit.ID
			err = a.db.Habits.CreateHabitLog(dataLog)
		} else {
			// if it exists update it
			hl.Origin = dataLog.Origin
			hl.Note = dataLog.Note
			hl.Success = dataLog.Success
			hl.IsAutomated = dataLog.IsAutomated
			_, err = a.db.Habits.UpdateHabitLog(hl)
		}
		if err != nil {
			return nil, err
		}
	}

	if err := database.ExecuteHabitsPipeline([]*entity.Habit{habit}, database.HabitsJoinLogs(a.db)); err != nil {
		return nil, err
	}
	habit.Status = calculateHabitStatusFromLogs(habit.Logs)

	return a.db.Habits.Update(habit, database.HabitsPipeline(a.db)...)
}

func (a *App) HabitFullUpdate(habit, data *entity.Habit) (*entity.Habit, error) {
	data.ID = habit.ID
	return a.db.Habits.Update(data, database.HabitsPipeline(a.db)...)
}

func (a *App) HabitDelete(habit *entity.Habit) error {
	return a.db.Habits.Delete(habit.ID)
}

func (a *App) HabitsGetAll(from, to time.Time) ([]*entity.Habit, error) {
	return a.db.Habits.Find(entity.HabitFilter{}, database.HabitsPipeline(a.db)...)
}

func (a *App) HabitGetByID(id int) (*entity.Habit, error) {
	return a.db.Habits.Get(entity.HabitFilter{ID: []int{id}}, database.HabitsPipeline(a.db)...)
}

func (a *App) HabitsGetCategories() ([]*entity.HabitCategory, error) {
	return a.db.Habits.FindHabitCategories(entity.HabitCategoryFilter{})
}

func (a *App) getOrCreateHabit(dayID, categoryID int) (*entity.Habit, error) {
	// First check that the habit already exists
	h, err := a.db.Habits.Get(entity.HabitFilter{DayID: []int{dayID}, CategoryID: []int{categoryID}})
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return nil, err
	}
	if h == nil {
		// If it does not exist we create it
		h = &entity.Habit{
			DayID:      dayID,
			CategoryID: categoryID,
			Status:     entity.HabitStatusNoInfo,
		}
		err = a.db.Habits.Create(h)
	}
	return h, err
}

func calculateHabitStatusFromLogs(logs []*entity.HabitLog) entity.HabitStatus {
	status := entity.HabitStatusNoInfo
	for _, log := range logs {
		if !log.IsAutomated {
			if log.Success {
				return entity.HabitStatusDone
			}
			status = entity.HabitStatusNotDone
		}
	}
	if status == entity.HabitStatusNoInfo {
		for _, log := range logs {
			if log.IsAutomated {
				if log.Success {
					return entity.HabitStatusDone
				}
				status = entity.HabitStatusNotDone
			}
		}
	}

	return status
}
