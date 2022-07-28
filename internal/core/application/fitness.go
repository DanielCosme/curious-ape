package application

import (
	"fmt"
	"github.com/danielcosme/curious-ape/internal/core/database"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/sdk/dates"
	"github.com/danielcosme/curious-ape/sdk/errors"
	"github.com/danielcosme/curious-ape/sdk/log"
	"time"
)

func (a *App) GetFitnessLogs(filter entity.FitnessLogFilter) ([]*entity.FitnessLog, error) {
	return a.db.FitnessLogs.Find(filter, database.FitnessLogsJoinDay(a.db))
}

func (a *App) GetFitnessLog(filter entity.FitnessLogFilter) (*entity.FitnessLog, error) {
	return a.db.FitnessLogs.Get(filter, database.FitnessLogsJoinDay(a.db))
}

func (a *App) DeleteFitnessLog(fl *entity.FitnessLog) error {
	return a.db.FitnessLogs.Delete(fl.ID)
}

func (a *App) SyncFitnessLog(date time.Time) error {
	client, err := a.Oauth2GetClient(entity.ProviderGoogle)
	if err != nil {
		return err
	}

	googleAPI := a.Sync.GoogleClient(client)
	_, err = googleAPI.Fitness.GetFitnessSessions(dates.ToBeginningOfDay(date), dates.ToEndOfDay(date))
	// TODO continue here
	return err
}

func (a *App) CreateFitnessLogFromApi(fitnessLog *entity.FitnessLog) (*entity.FitnessLog, error) {
	if fitnessLog.Raw == "" {
		// TODO
	}

	if err := a.createFitnessLogs([]*entity.FitnessLog{fitnessLog}); err != nil {
		return nil, err
	}
	return a.db.FitnessLogs.Get(entity.FitnessLogFilter{ID: []int{fitnessLog.ID}}, database.FitnessLogsJoinDay(a.db))
}

func (a *App) UpdateFitnessLog(fl, data *entity.FitnessLog) (*entity.FitnessLog, error) {
	if err := checkManualOrigin(fl.Origin); err != nil {
		return nil, err
	}

	data.ID = fl.ID
	return a.db.FitnessLogs.Update(data, database.FitnessLogsJoinDay(a.db))
}

func (a *App) createFitnessLogs(fls []*entity.FitnessLog) error {
	for _, fl := range fls {
		if err := a.HabitUpsertFromFitnessLog(fl); err != nil {
			return err
		}

		fl.Date = fl.Day.Date
		if err := a.db.FitnessLogs.Create(fl); err != nil {
			return err
		}

		a.Log.InfoP("Created fitness log", log.Prop{
			"provider": fl.Origin.Str(),
			"date":     fl.Day.Date.Format(entity.HumanDate),
			"type":     fl.Type.Str(),
		})
	}
	return nil
}

func checkManualOrigin(o entity.DataSource) error {
	if o != entity.Manual {
		return errors.New("only manually crated logs can be updated")

	}
	return nil
}

func (a *App) HabitUpsertFromFitnessLog(fl *entity.FitnessLog) error {
	habitCategory, err := a.GetHabitCategoryByType(entity.HabitTypeFitness)
	if err != nil {
		return err
	}

	var success bool
	if fl.Type == entity.StrengthTraining {
		success = true

		habit := &entity.Habit{
			DayID:      fl.DayID,
			CategoryID: habitCategory.ID,
			Logs: []*entity.HabitLog{{
				Success:     success,
				IsAutomated: fl.Origin == entity.Manual,
				Origin:      fl.Origin,
				Note:        fmt.Sprintf("Fitness log of type %s", fl.Type),
			}},
		}
		habit, err = a.HabitCreate(fl.Day, habit)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) GetHabitCategoryByType(t entity.HabitType) (*entity.HabitCategory, error) {
	return a.db.Habits.GetHabitCategory(entity.HabitCategoryFilter{Type: []entity.HabitType{t}})
}
