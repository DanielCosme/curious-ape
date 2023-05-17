package application

import (
	"fmt"
	"strconv"
	"time"

	db "github.com/danielcosme/curious-ape/internal/core/database"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/go-sdk/log"
)

type NewHabitParams struct {
	Date         time.Time
	CategoryCode string
	Success      bool
	Origin       entity.DataSource
	Note         string
	IsAutomated  bool
}

func (p *NewHabitParams) ToLog() *entity.HabitLog {
	return &entity.HabitLog{
		Success:     p.Success,
		Note:        p.Note,
		Origin:      p.Origin,
		IsAutomated: p.IsAutomated,
	}
}

func (a *App) HabitCreate(data *NewHabitParams) (*entity.Habit, error) {
	habit, err := db.GetOrCreateHabit(a.db, data.Date, data.CategoryCode, db.HabitsPipeline(a.db)...)
	if err != nil {
		return nil, err
	}

	hl := data.ToLog()
	hl.HabitID = habit.ID
	operation, err := db.UpsertHabitLog(a.db, hl)
	if err != nil {
		return nil, err
	}

	a.Log.InfoP(fmt.Sprintf("habit log succesfully %s", operation), log.Prop{
		"Type":    habit.Category.Type.Str(),
		"Success": strconv.FormatBool(hl.Success),
		"Origin":  hl.Origin.Str(),
		"details": hl.Note,
		"date":    entity.FormatDate(data.Date),
	})

	if err := db.ExecuteHabitsPipeline([]*entity.Habit{habit}, db.HabitsJoinLogs(a.db)); err != nil {
		return nil, err
	}

	oldStatus := habit.Status
	habit.CalculateHabitStatus()
	if oldStatus != habit.Status {
		return a.db.Habits.Update(habit, db.HabitsPipeline(a.db)...)
	}

	return habit, nil
}

func (a *App) HabitUpsertFromSleepLog(sleepLog entity.SleepLog) error {
	if !sleepLog.IsMainSleep {
		return nil
	}

	var success bool
	wakeUPTime := toWakeUpTime(sleepLog.EndTime)
	if sleepLog.EndTime.Before(wakeUPTime) {
		success = true
	}

	_, err := a.HabitCreate(&NewHabitParams{
		Date:         sleepLog.Date,
		CategoryCode: entity.HabitTypeWakeUp.Str(),
		Success:      success,
		Origin:       sleepLog.Origin,
		Note:         fmt.Sprintf("Wake up time %s", sleepLog.EndTime.Format(entity.Timestamp)),
		IsAutomated:  false,
	})
	if err != nil {
		return err
	}

	return nil
}

// 6 a.m.
func toWakeUpTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 6, 0, 0, 0, t.Location())
}

func (a *App) HabitFullUpdate(habit, data *entity.Habit) (*entity.Habit, error) {
	data.ID = habit.ID
	return a.db.Habits.Update(data, db.HabitsPipeline(a.db)...)
}

func (a *App) HabitDelete(habit *entity.Habit) error {
	return a.db.Habits.Delete(habit.ID)
}

func (a *App) HabitsGetAll(params map[string]string) ([]*entity.Habit, error) {
	f := entity.HabitFilter{}

	var err error
	end, start := time.Time{}, time.Time{}
	if len(params) > 0 {
		if d, ok := params["startDate"]; ok {
			start, err = entity.ParseDate(d)
			if err != nil {
				return nil, err
			}
		}
		if d, ok := params["endDate"]; ok {
			end, err = entity.ParseDate(d)
			if err != nil {
				return nil, err
			}
		}

		days, err := a.daysGetByDateRange(start, end)
		if err != nil {
			return nil, err
		}
		f.DayID = db.DayToIDs(days)
	}

	return a.db.Habits.Find(f, db.HabitsPipeline(a.db)...)
}

func (a *App) HabitGetByID(id int) (*entity.Habit, error) {
	return a.db.Habits.Get(entity.HabitFilter{ID: []int{id}}, db.HabitsPipeline(a.db)...)
}

func (a *App) HabitsGetByDay(d *entity.Day) ([]*entity.Habit, error) {
	return a.db.Habits.Find(entity.HabitFilter{DayID: []int{d.ID}}, db.HabitsPipeline(a.db)...)
}

func (a *App) HabitsGetCategories() ([]*entity.HabitCategory, error) {
	return a.db.Habits.FindHabitCategories(entity.HabitCategoryFilter{})
}
