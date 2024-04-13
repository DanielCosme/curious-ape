package application

import (
	"fmt"
	database2 "github.com/danielcosme/curious-ape/internal/database"
	entity2 "github.com/danielcosme/curious-ape/internal/entity"
	"strconv"
	"time"

	"github.com/danielcosme/go-sdk/log"
)

type NewHabitParams struct {
	Date         time.Time
	CategoryCode string
	Success      bool
	Origin       entity2.DataSource
	Note         string
	IsAutomated  bool
}

func (p *NewHabitParams) ToLog() *entity2.HabitLog {
	return &entity2.HabitLog{
		Success:     p.Success,
		Note:        p.Note,
		Origin:      p.Origin,
		IsAutomated: p.IsAutomated,
	}
}

func (a *App) HabitUpsert(params *NewHabitParams) (*entity2.Habit, error) {
	habit, err := database2.GetOrCreateHabit(a.db, params.Date, params.CategoryCode)
	if err != nil {
		return nil, err
	}

	hl := params.ToLog()
	hl.HabitID = habit.ID
	operation, err := database2.UpsertHabitLog(a.db, hl)
	if err != nil {
		return nil, err
	}
	if err := database2.ExecuteHabitsPipeline([]*entity2.Habit{habit}, database2.HabitsPipeline(a.db)...); err != nil {
		return nil, err
	}

	a.Log.InfoP(fmt.Sprintf("habit log succesfully %s", operation), log.Prop{
		"Type":    habit.Category.Type.Str(),
		"Success": strconv.FormatBool(hl.Success),
		"Origin":  hl.Origin.Str(),
		"Details": hl.Note,
		"Date":    entity2.FormatDate(params.Date),
	})

	oldStatus := habit.Status
	habit.Status = entity2.CalculateHabitStatus(habit.Logs)
	if oldStatus != habit.Status {
		a.Log.InfoP("habit status changed", log.Prop{
			"From": string(oldStatus),
			"To":   string(habit.Status),
			"Date": entity2.FormatDate(params.Date),
			"Type": habit.Category.Type.Str(),
		})
		return a.db.Habits.Update(habit, database2.HabitsPipeline(a.db)...)
	}
	return habit, nil
}

func (a *App) HabitUpsertFromSleepLog(sleepLog entity2.SleepLog) error {
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
		CategoryCode: entity2.HabitTypeWakeUp.Str(),
		Success:      success,
		Origin:       sleepLog.Origin,
		Note:         fmt.Sprintf("Wake up time %s", sleepLog.EndTime.Format(entity2.Timestamp)),
		IsAutomated:  true,
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

func (a *App) HabitFullUpdate(habit, data *entity2.Habit) (*entity2.Habit, error) {
	data.ID = habit.ID
	return a.db.Habits.Update(data, database2.HabitsPipeline(a.db)...)
}

func (a *App) HabitDelete(habit *entity2.Habit) error {
	return a.db.Habits.Delete(habit.ID)
}

func (a *App) HabitsGetAll(params map[string]string) ([]*entity2.Habit, error) {
	f := entity2.HabitFilter{}

	var err error
	end, start := time.Time{}, time.Time{}
	if len(params) > 0 {
		if d, ok := params["startDate"]; ok {
			start, err = entity2.ParseDate(d)
			if err != nil {
				return nil, err
			}
		}
		if d, ok := params["endDate"]; ok {
			end, err = entity2.ParseDate(d)
			if err != nil {
				return nil, err
			}
		}

		days, err := a.daysGetByDateRange(start, end)
		if err != nil {
			return nil, err
		}
		f.DayID = database2.DayToIDs(days)
	}

	return a.db.Habits.Find(f, database2.HabitsPipeline(a.db)...)
}

func (a *App) HabitGetByID(id int) (*entity2.Habit, error) {
	return a.db.Habits.Get(entity2.HabitFilter{ID: []int{id}}, database2.HabitsPipeline(a.db)...)
}

func (a *App) HabitsGetByDay(d *entity2.Day) ([]*entity2.Habit, error) {
	return a.db.Habits.Find(entity2.HabitFilter{DayID: []int{d.ID}}, database2.HabitsPipeline(a.db)...)
}

func (a *App) HabitsGetCategories() ([]*entity2.HabitCategory, error) {
	return a.db.Habits.FindHabitCategories(entity2.HabitCategoryFilter{})
}
