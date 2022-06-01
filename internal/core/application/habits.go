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

	h, err := a.getOrCreateHabit(day.ID, habitCategory.ID)
	if err != nil {
		return nil, err
	}

	// TODO hardcode origins for habit logs and create a validator for it
	// Create the habit log
	for _, l := range data.Logs {
		hl, err := a.db.Habits.GetHabitLog(entity.HabitFilter{Origin: l.Origin, ID: h.ID})
		if err != nil {
			a.Log.Trace("Habit log err", err.Error())
			if errors.Is(err, entity.ErrNotFound) {
				// if it does not exist create it
				l.HabitID = h.ID
				err := a.db.Habits.CreateHabitLog(l)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			a.Log.Trace("Habit log exists then we update")
			// if it exists update it
			hl.Origin = l.Origin
			hl.Note = l.Note
			hl.Success = l.Success
			hl.IsAutomated = l.IsAutomated
			hl, err = a.db.Habits.UpdateHabitLog(hl)
			if err != nil {
				return nil, err
			}
		}

		// Calculate status
		if l.Success {
			h.Status = entity.HabitStatusDone
		} else {
			h.Status = entity.HabitStatusNotDone
		}
		_, err = a.db.Habits.Update(h)
		if err != nil {
			return nil, err
		}
	}

	return h, repository.ExecuteHabitsPipeline([]*entity.Habit{h}, datasource.HabitsPipeline(a.db)...)
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
