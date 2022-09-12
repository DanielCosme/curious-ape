package application

import (
	"encoding/json"
	"fmt"
	"github.com/danielcosme/curious-ape/internal/core/database"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/integrations/google"
	"github.com/danielcosme/go-sdk/dates"
	"github.com/danielcosme/go-sdk/errors"
	"github.com/danielcosme/go-sdk/log"
	"strconv"
	"time"
)

func (a *App) FitnessFindLogs(filter entity.FitnessLogFilter) ([]*entity.FitnessLog, error) {
	return a.db.FitnessLogs.Find(filter, database.FitnessLogsJoinDay(a.db))
}

func (a *App) FitnessGetLog(filter entity.FitnessLogFilter) (*entity.FitnessLog, error) {
	return a.db.FitnessLogs.Get(filter, database.FitnessLogsJoinDay(a.db))
}

func (a *App) FitnessDeleteLog(fl *entity.FitnessLog) error {
	return a.db.FitnessLogs.Delete(fl.ID)
}

func (a *App) SyncFitness() error {
	days, err := a.db.Days.Find(entity.DayFilter{}, database.DaysJoinFitnessLogs(a.db))
	if err != nil {
		return err
	}
	client, err := a.Oauth2GetClient(entity.ProviderGoogle)
	if err != nil {
		return err
	}
	googleApi := a.sync.GoogleClient(client)

	for _, d := range days {
		// TODO this will break for multiple providers and types of fitness logs (steps, cycling, etc...)
		if len(d.FitnessLogs) == 0 {
			payload, err := googleApi.Fitness.GetFitnessSessions(dates.BeginningOfDay(d.Date), dates.EndOfDay(d.Date))
			if err != nil {
				return err
			}
			fitnessLogs, withoutLog, err := toFitnessLogFromGoogle([]*entity.Day{d}, payload)
			if err != nil {
				return err
			}

			a.createFailedHabitForDays(withoutLog, entity.HabitTypeFitness, entity.Google)
			if err := a.createFitnessLogs(fitnessLogs); err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *App) SyncFitnessByDateRAnge(start, end time.Time) error {
	client, err := a.Oauth2GetClient(entity.ProviderGoogle)
	if err != nil {
		return err
	}
	days, err := a.daysGetByDateRange(start, end)
	if err != nil {
		return err
	}

	googleAPI := a.sync.GoogleClient(client)
	gFitnessLogs, err := googleAPI.Fitness.GetFitnessSessions(dates.BeginningOfDay(start), dates.EndOfDay(end))
	if err != nil {
		return err
	}

	fitnessLogs, daysWithoutLog, err := toFitnessLogFromGoogle(days, gFitnessLogs)
	if err != nil {
		return err
	}
	// In case of missing logs for day create the habit log in failed state.
	a.createFailedHabitForDays(daysWithoutLog, entity.HabitTypeFitness, entity.Google)

	return a.createFitnessLogs(fitnessLogs)
}

func (a *App) SyncFitnessLog(date time.Time) error {
	day, err := a.DayGetByDate(date)
	if err != nil {
		return err
	}

	client, err := a.Oauth2GetClient(entity.ProviderGoogle)
	if err != nil {
		return err
	}
	googleAPI := a.sync.GoogleClient(client)
	gFitnessLogs, err := googleAPI.Fitness.GetFitnessSessions(dates.BeginningOfDay(date), dates.EndOfDay(date))
	if err != nil {
		return err
	}

	fitnessLogs, daysWithoutLog, err := toFitnessLogFromGoogle([]*entity.Day{day}, gFitnessLogs)
	if err != nil {
		return err
	}
	// In case of missing logs for day create the habit log in failed state.
	a.createFailedHabitForDays(daysWithoutLog, entity.HabitTypeFitness, entity.Google)

	return a.createFitnessLogs(fitnessLogs)
}

func toFitnessLogFromGoogle(days []*entity.Day, gfls []google.Session) ([]*entity.FitnessLog, []*entity.Day, error) {
	fitnessLogs := make([]*entity.FitnessLog, 0, len(gfls))
	daysWithWithoutLog := []*entity.Day{}

	mapGoogleFitnessLogByDate := map[string]google.Session{}
	for _, gfl := range gfls {
		millis, _ := strconv.Atoi(gfl.StartTimeMillis)
		dateOfWorkout := time.UnixMilli(int64(millis))
		mapGoogleFitnessLogByDate[dateOfWorkout.Format(entity.ISO8601)] = gfl
	}

	for _, day := range days {
		if gfl, ok := mapGoogleFitnessLogByDate[day.FormatDate()]; ok {
			millis, _ := strconv.Atoi(gfl.StartTimeMillis)
			startTime := time.UnixMilli(int64(millis))
			millisEnd, _ := strconv.Atoi(gfl.EndTimeMillis)
			endTime := time.UnixMilli(int64(millisEnd))

			fitnessLog := &entity.FitnessLog{
				DayID:     day.ID,
				Title:     gfl.Name,
				Date:      day.Date,
				StartTime: startTime,
				EndTime:   endTime,
				Origin:    entity.Google,
				Note:      fmt.Sprintf("from: %s", gfl.Application.PackageName),
				Day:       day,
			}

			switch gfl.ActivityType {
			case 97: // weightlifting
				fitnessLog.Type = entity.StrengthTraining
			default:
				return nil, nil, errors.New("google: cannot determine type of activity type")
			}

			raw, err := json.Marshal(gfl)
			if err != nil {
				return nil, nil, err
			}

			fitnessLog.Raw = string(raw)
			fitnessLogs = append(fitnessLogs, fitnessLog)
		} else {
			daysWithWithoutLog = append(daysWithWithoutLog, day)
		}
	}

	return fitnessLogs, daysWithWithoutLog, nil
}

func (a *App) createFailedHabitForDays(days []*entity.Day, category entity.HabitType, source entity.DataSource) {
	habitCategory, err := a.HabitCategoryGetByType(category)
	if err != nil {
		a.Log.Error(err)
	}

	for _, day := range days {
		habit := &entity.Habit{
			DayID:      day.ID,
			CategoryID: habitCategory.ID,
			Logs: []*entity.HabitLog{{
				Success:     false,
				IsAutomated: source != entity.Manual,
				Origin:      source,
				Note:        fmt.Sprintf("From missing log on data source"),
			}},
		}
		habit, err = a.HabitCreate(day, habit)
		if err != nil {
			a.Log.Error(err)
		}
	}
}

func (a *App) FitnessCreateLogFromApi(fitnessLog *entity.FitnessLog) (*entity.FitnessLog, error) {
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
		fl.Date = fl.Day.Date
		if err := a.db.FitnessLogs.Create(fl); err != nil {
			a.Log.Warningf("fitness log for %s could not be created: %s", fl.Day.Date.Format(entity.HumanDate), err.Error())
			continue
		}

		if err := a.HabitUpsertFromFitnessLog(fl); err != nil {
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
	habitCategory, err := a.HabitCategoryGetByType(entity.HabitTypeFitness)
	if err != nil {
		return err
	}

	habit := &entity.Habit{
		DayID:      fl.DayID,
		CategoryID: habitCategory.ID,
		Logs: []*entity.HabitLog{{
			Success:     true,
			IsAutomated: fl.Origin != entity.Manual,
			Origin:      fl.Origin,
			Note:        fmt.Sprintf("Fitness log of type %s", fl.Type),
		}},
	}
	habit, err = a.HabitCreate(fl.Day, habit)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) HabitCategoryGetByType(t entity.HabitType) (*entity.HabitCategory, error) {
	return a.db.Habits.GetHabitCategory(entity.HabitCategoryFilter{Type: []entity.HabitType{t}})
}
