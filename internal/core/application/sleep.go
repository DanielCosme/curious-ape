package application

import (
	fitbit2 "github.com/danielcosme/curious-ape/fitbit/fitbit"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"time"
)

func (a *App) SleepDebug() (*fitbit2.SleepEnvelope, error) {
	// Get client, refreshes token if necessary
	client, err := a.Oauth2GetClient(entity.ProviderFitbit)
	if err != nil {
		return nil, err
	}

	api := fitbit2.NewAPI(client)
	sleepLog, err := api.Sleep.GetLogByDate(time.Now().AddDate(0, 0, 0))

	return sleepLog, err
}
