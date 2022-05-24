package application

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/integrations/fitbit"
	"time"
)

func (a *App) SleepDebug() (*fitbit.SleepEnvelope, error) {
	// Get client, refreshes token if necessary
	client, err := a.Oauth2GetClient(entity.ProviderFitbit)
	if err != nil {
		return nil, err
	}

	api := fitbit.NewAPI(client)
	sleepLog, err := api.Sleep.GetLogByDate(time.Now())

	return sleepLog, err
}
