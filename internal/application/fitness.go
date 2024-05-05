package application

import (
	"encoding/json"
	"errors"
	"fmt"
	database2 "github.com/danielcosme/curious-ape/internal/database"
	"github.com/danielcosme/curious-ape/internal/dates"
	entity2 "github.com/danielcosme/curious-ape/internal/entity"
	"strconv"
	"time"

	"github.com/danielcosme/curious-ape/internal/integrations/google"
)

func (a *App) FitnessFindLogs(filter entity2.FitnessLogFilter) ([]*entity2.FitnessLog, error) {
	return a.db.FitnessLogs.Find(filter, database2.FitnessLogsJoinDay(a.db))
}

func (a *App) FitnessGetLog(filter entity2.FitnessLogFilter) (*entity2.FitnessLog, error) {
	return a.db.FitnessLogs.Get(filter, database2.FitnessLogsJoinDay(a.db))
}

func (a *App) FitnessDeleteLog(fl *entity2.FitnessLog) error {
	return a.db.FitnessLogs.Delete(fl.ID)
}

func (a *App) SyncFitness() error {
	days, err := a.db.Days.Find(entity2.DayFilter{}, database2.DaysJoinFitnessLogs(a.db))
	if err != nil {
		return err
	}
	client, err := a.Oauth2GetClient(entity2.ProviderGoogle)
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
			fitnessLogs, withoutLog, err := toFitnessLogFromGoogle([]*entity2.Day{d}, payload)
			if err != nil {
				return err
			}

			a.createFailedHabitForDays(withoutLog, entity2.HabitTypeFitness, entity2.Google)
			if err := a.createFitnessLogs(fitnessLogs); err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *App) SyncFitnessByDateRAnge(start, end time.Time) error {
	client, err := a.Oauth2GetClient(entity2.ProviderGoogle)
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
	a.createFailedHabitForDays(daysWithoutLog, entity2.HabitTypeFitness, entity2.Google)

	return a.createFitnessLogs(fitnessLogs)
}

func (a *App) SyncFitnessLog(date time.Time) error {
	day, err := database2.DayGetOrCreate(a.db, date)
	if err != nil {
		return err
	}

	client, err := a.Oauth2GetClient(entity2.ProviderGoogle)
	if err != nil {
		return err
	}
	googleAPI := a.sync.GoogleClient(client)
	gFitnessLogs, err := googleAPI.Fitness.GetFitnessSessions(dates.BeginningOfDay(date), dates.EndOfDay(date))
	if err != nil {
		return err
	}

	fitnessLogs, daysWithoutLog, err := toFitnessLogFromGoogle([]*entity2.Day{day}, gFitnessLogs)
	if err != nil {
		return err
	}
	// In case of missing logs for day create the habit log in failed state.
	a.createFailedHabitForDays(daysWithoutLog, entity2.HabitTypeFitness, entity2.Google)

	return a.createFitnessLogs(fitnessLogs)
}

func toFitnessLogFromGoogle(days []*entity2.Day, gfls []google.Session) ([]*entity2.FitnessLog, []*entity2.Day, error) {
	fitnessLogs := make([]*entity2.FitnessLog, 0, len(gfls))
	daysWithWithoutLog := []*entity2.Day{}

	mapGoogleFitnessLogByDate := map[string]google.Session{}
	for _, gfl := range gfls {
		millis, _ := strconv.Atoi(gfl.StartTimeMillis)
		dateOfWorkout := time.UnixMilli(int64(millis))
		mapGoogleFitnessLogByDate[dateOfWorkout.Format(entity2.ISO8601)] = gfl
	}

	for _, day := range days {
		if gfl, ok := mapGoogleFitnessLogByDate[day.FormatDate()]; ok {
			millis, _ := strconv.Atoi(gfl.StartTimeMillis)
			startTime := time.UnixMilli(int64(millis))
			millisEnd, _ := strconv.Atoi(gfl.EndTimeMillis)
			endTime := time.UnixMilli(int64(millisEnd))

			fitnessLog := &entity2.FitnessLog{
				DayID:     day.ID,
				Title:     gfl.Name,
				Date:      day.Date,
				StartTime: startTime,
				EndTime:   endTime,
				Origin:    entity2.Google,
				Note:      fmt.Sprintf("from: %s", gfl.Application.PackageName),
				Day:       day,
			}

			switch gfl.ActivityType {
			case 97: // weightlifting
				fitnessLog.Type = entity2.StrengthTraining
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

func (a *App) createFailedHabitForDays(days []*entity2.Day, category entity2.HabitType, source entity2.DataSource) {
	habitCategory, err := a.HabitCategoryGetByType(category)
	if err != nil {
		a.Log.Error(err.Error())
	}

	for _, day := range days {
		_, err = a.HabitUpsert(&NewHabitParams{
			Date:         day.Date,
			CategoryCode: habitCategory.Code,
			Success:      false,
			IsAutomated:  source != entity2.Manual,
			Origin:       source,
			Note:         fmt.Sprintf("From missing log on data source"),
		})
		if err != nil {
			a.Log.Error(err.Error())
		}
	}
}

func (a *App) FitnessCreateLogFromApi(fitnessLog *entity2.FitnessLog) (*entity2.FitnessLog, error) {
	if fitnessLog.Raw == "" {
		// TODO
	}

	if err := a.createFitnessLogs([]*entity2.FitnessLog{fitnessLog}); err != nil {
		return nil, err
	}
	return a.db.FitnessLogs.Get(entity2.FitnessLogFilter{ID: []int{fitnessLog.ID}}, database2.FitnessLogsJoinDay(a.db))
}

func (a *App) UpdateFitnessLog(fl, data *entity2.FitnessLog) (*entity2.FitnessLog, error) {
	if err := checkManualOrigin(fl.Origin); err != nil {
		return nil, err
	}

	data.ID = fl.ID
	return a.db.FitnessLogs.Update(data, database2.FitnessLogsJoinDay(a.db))
}

func (a *App) createFitnessLogs(fls []*entity2.FitnessLog) error {
	for _, fl := range fls {
		fl.Date = fl.Day.Date
		if err := a.db.FitnessLogs.Create(fl); err != nil {
			a.Log.Warn("fitness log for %s could not be created: %s", fl.Day.Date.Format(entity2.HumanDateWithTime), err.Error())
			continue
		}

		if err := a.HabitUpsertFromFitnessLog(fl); err != nil {
			return err
		}

		a.Log.Info("Created fitness log",
			"provider", fl.Origin.Str(),
			"date", fl.Day.Date.Format(entity2.HumanDateWithTime),
			"type", fl.Type.Str(),
		)
	}
	return nil
}

func checkManualOrigin(o entity2.DataSource) error {
	if o != entity2.Manual {
		return errors.New("only manually crated logs can be updated")
	}
	return nil
}

func (a *App) HabitUpsertFromFitnessLog(fl *entity2.FitnessLog) error {
	habitCategory, err := a.HabitCategoryGetByType(entity2.HabitTypeFitness)
	if err != nil {
		return err
	}

	_, err = a.HabitUpsert(&NewHabitParams{
		Date:         fl.Date,
		CategoryCode: habitCategory.Code,
		Success:      true,
		IsAutomated:  fl.Origin != entity2.Manual,
		Origin:       fl.Origin,
		Note:         fmt.Sprintf("Fitness log of type %s", fl.Type),
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *App) HabitCategoryGetByType(t entity2.HabitType) (*entity2.HabitCategory, error) {
	return a.db.Habits.GetHabitCategory(entity2.HabitCategoryFilter{Type: []entity2.HabitType{t}})
}
