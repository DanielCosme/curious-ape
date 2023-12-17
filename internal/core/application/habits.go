package application

import (
	"fmt"
	"strconv"
	"time"

	"github.com/danielcosme/curious-ape/internal/core/database"
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

func (a *App) HabitUpsert(params *NewHabitParams) (*entity.Habit, error) {
	habit, err := database.GetOrCreateHabit(a.db, params.Date, params.CategoryCode)
	if err != nil {
		return nil, err
	}

	hl := params.ToLog()
	hl.HabitID = habit.ID
	operation, err := database.UpsertHabitLog(a.db, hl)
	if err != nil {
		return nil, err
	}
	if err := database.ExecuteHabitsPipeline([]*entity.Habit{habit}, database.HabitsPipeline(a.db)...); err != nil {
		return nil, err
	}

	a.Log.InfoP(fmt.Sprintf("habit log succesfully %s", operation), log.Prop{
		"Type":    habit.Category.Type.Str(),
		"Success": strconv.FormatBool(hl.Success),
		"Origin":  hl.Origin.Str(),
		"Details": hl.Note,
		"Date":    entity.FormatDate(params.Date),
	})

	oldStatus := habit.Status
	habit.Status = entity.CalculateHabitStatus(habit.Logs)
	if oldStatus != habit.Status {
		a.Log.InfoP("habit status changed", log.Prop{
			"From": string(oldStatus),
			"To":   string(habit.Status),
			"Date": entity.FormatDate(params.Date),
			"Type": habit.Category.Type.Str(),
		})
		return a.db.Habits.Update(habit, database.HabitsPipeline(a.db)...)
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

	_, err := a.HabitUpsert(&NewHabitParams{
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
	return a.db.Habits.Update(data, database.HabitsPipeline(a.db)...)
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
		f.DayID = database.DayToIDs(days)
	}

	return a.db.Habits.Find(f, database.HabitsPipeline(a.db)...)
}

func (a *App) HabitGetByID(id int) (*entity.Habit, error) {
	return a.db.Habits.Get(entity.HabitFilter{ID: []int{id}}, database.HabitsPipeline(a.db)...)
}

func (a *App) HabitsGetByDay(d *entity.Day) ([]*entity.Habit, error) {
	return a.db.Habits.Find(entity.HabitFilter{DayID: []int{d.ID}}, database.HabitsPipeline(a.db)...)
}

func (a *App) HabitsGetCategories() ([]*entity.HabitCategory, error) {
	return a.db.Habits.FindHabitCategories(entity.HabitCategoryFilter{})
}
