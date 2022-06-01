package application

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/core/repository"
	"github.com/danielcosme/curious-ape/internal/datasource"
	"github.com/danielcosme/curious-ape/sdk/errors"
	"time"
)

func (a *App) HabitCreate(day *entity.Day, data *entity.Habit) (*entity.Habit, error) {
	habitCategory, err := a.db.Habits.GetHabitCategory(entity.HabitFilter{CategoryID: data.CategoryID})
	if err != nil {
		return nil, err
	}

	habit, err := a.getOrCreateHabit(day.ID, habitCategory.ID)
	if err != nil {
		return nil, err
	}

	// TODO hardcode origins for habit logs and create a validator for it
	// Create the habit log
	for _, l := range data.Logs {
		hl, err := a.db.Habits.GetHabitLog(entity.HabitFilter{Origin: l.Origin, ID: habit.ID})
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {
				// if it does not exist create it
				l.HabitID = habit.ID
				err := a.db.Habits.CreateHabitLog(l)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			// if it exists update it
			hl.Origin = l.Origin
			hl.Note = l.Note
			hl.Success = l.Success
			hl.IsAutomated = l.IsAutomated
			_, err = a.db.Habits.UpdateHabitLog(hl)
			if err != nil {
				return nil, err
			}
		}
	}

	// Calculate habit status based on the logs
	if err := repository.ExecuteHabitsPipeline([]*entity.Habit{habit}, datasource.HabitsJoinLogs(a.db)); err != nil {
		return nil, err
	}

	status := entity.HabitStatusNoInfo
	// re-calculate habit status
	for _, log := range habit.Logs {
		if !log.IsAutomated {
			// if the log is manually added
			if log.Success {
				status = entity.HabitStatusDone
			} else {
				status = entity.HabitStatusNotDone
			}
			break
		} else {
			if log.Success {
				status = entity.HabitStatusDone
			} else if status == entity.HabitStatusNoInfo {
				status = entity.HabitStatusNotDone
			}
		}
	}
	habit.Status = status

	return a.db.Habits.Update(habit, datasource.HabitsPipeline(a.db)...)
}

func (a *App) HabitFullUpdate(habit, data *entity.Habit) (*entity.Habit, error) {

	data.ID = habit.ID
	return a.db.Habits.Update(data, datasource.HabitsPipeline(a.db)...)
}

func (a *App) HabitDelete(habit *entity.Habit) error {
	return a.db.Habits.Delete(habit.ID)
}

func (a *App) HabitsGetAll(from, to time.Time) ([]*entity.Habit, error) {
	return a.db.Habits.Find(entity.HabitFilter{}, datasource.HabitsPipeline(a.db)...)
}

func (a *App) HabitGetByID(id int) (*entity.Habit, error) {
	return a.db.Habits.Get(entity.HabitFilter{ID: id}, datasource.HabitsPipeline(a.db)...)
}

func (a *App) HabitsGetCategories() ([]*entity.HabitCategory, error) {
	return a.db.Habits.FindHabitCategories(entity.HabitFilter{})
}

func (a *App) getOrCreateHabit(dayID, categoryID int) (*entity.Habit, error) {
	// First check that the habit already exists
	h, err := a.db.Habits.Get(entity.HabitFilter{DayID: dayID, CategoryID: categoryID})
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			// If it does not exist we create it
			h = &entity.Habit{
				DayID:      dayID,
				CategoryID: categoryID,
				Status:     entity.HabitStatusNoInfo,
			}
			if err = a.db.Habits.Create(h); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return h, nil
}
