package application

import (
	"github.com/danielcosme/curious-ape/internal/core/database"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/sdk/dates"
	"time"
)

func (a *App) GetFitnessLogsForDay(d *entity.Day) ([]*entity.FitnessLog, error) {
	return a.getFitnessLogs(entity.FitnessLogFilter{DayID: []int{d.ID}})
}

func (a *App) getFitnessLogs(filter entity.FitnessLogFilter) ([]*entity.FitnessLog, error) {
	return a.db.FitnessLogs.Find(filter, database.FitnessLogsJoinDay(a.db))
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
