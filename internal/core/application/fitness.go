package application

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/sdk/dates"
	"time"
)

func (a *App) FitnessGet() error {
	client, err := a.Oauth2GetClient(entity.ProviderGoogle)
	if err != nil {
		return err
	}

	googleAPI := a.Sync.GoogleClient(client)
	start := time.Now()
	_, err = googleAPI.Fitness.GetFitnessSessions(dates.ToBeginningOfDay(start), dates.ToEndOfDay(start))
	return err
}
