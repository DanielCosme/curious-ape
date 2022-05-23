package application

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/datasource"
	"time"
)

func (a *App) HabitCreate(d *entity.Day, h *entity.Habit) (*entity.Habit, error) {
	hc, err := a.db.Habits.GetHabitCategory(entity.HabitFilter{ID: []int{h.CategoryID}})
	if err != nil {
		return nil, err
	}

	if h.Origin == "" {
		h.Origin = entity.HabitOriginUnknown
	}

	h.DayID = d.ID
	h.CategoryID = hc.ID
	if err := a.db.Habits.Create(h, datasource.HabitsPipeline(a.db)...); err != nil {
		return nil, err
	}

	return h, nil
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
	return a.db.Habits.Get(entity.HabitFilter{ID: []int{id}}, datasource.HabitsPipeline(a.db)...)
}

func (a *App) HabitsGetCategories() ([]*entity.HabitCategory, error) {
	return a.db.Habits.FindHabitCategories(entity.HabitFilter{})
}
