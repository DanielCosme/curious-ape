package application

import (
	"fmt"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database"
	"github.com/danielcosme/curious-ape/internal/database/gen/models"
	"github.com/danielcosme/curious-ape/internal/integrations/fitbit"
	"golang.org/x/oauth2"
	"net/http"
)

type IntegrationInfo struct {
	Name        string
	State       IntegrationState
	ProfileInfo string
	AuthURL     string
	Problem     string
}

func (a *App) IntegrationsGet() ([]IntegrationInfo, error) {
	var res []IntegrationInfo
	currentIntegrations := []core.Integration{core.IntegrationFitbit}

	for _, integration := range currentIntegrations {
		var profileInfo, authURL, problem string
		state := IntegrationDisconnected

		switch integration {
		case core.IntegrationFitbit:
			sls, err := a.sleepLogsGetFromFitbit(core.NewDateToday())
			if err != nil {
				authURL = a.sync.GenerateOauth2URI(integration)
				problem = err.Error()
			} else {
				state = IntegrationConnected
				if len(sls) > 0 {
					profileInfo = fmt.Sprintf("Total time asleep last night: %s", sls[0].MinutesAsleep)
				}
			}
			res = append(res, IntegrationInfo{
				Name:        "Fitbit",
				State:       state,
				ProfileInfo: profileInfo,
				AuthURL:     authURL,
				Problem:     problem,
			})
		default:
		}
	}

	return res, nil
}

func (a *App) Oauth2Success(provider, code string) error {
	token, err := a.sync.ExchangeToken(core.Integration(provider), code)
	if err != nil {
		return err
	}
	_, err = a.db.Auths.Upsert(&models.AuthSetter{
		Provider:     omit.From(provider),
		AccessToken:  omit.From(token.AccessToken),
		RefreshToken: omitnull.From(token.RefreshToken),
		TokenType:    omitnull.From(token.Type()),
		Expiration:   omitnull.From(token.Expiry),
	})
	a.Log.Info("Authentication successful", "provider", provider, "code", code)
	return err
}

func (a *App) fitbitClient() (res fitbit.API, err error) {
	client, err := a.integrationsGetClient(core.IntegrationFitbit)
	if err != nil {
		return res, err
	}
	res = fitbit.NewAPI(client)
	return
}

func (a *App) integrationsGetClient(integration core.Integration) (*http.Client, error) {
	o, err := a.db.Auths.Get(database.AuthParams{Provider: integration})
	if err != nil {
		return nil, err
	}
	currentToken := &oauth2.Token{
		AccessToken:  o.AccessToken,
		RefreshToken: o.RefreshToken.GetOrZero(),
		Expiry:       o.Expiration.GetOrZero(),
		TokenType:    o.TokenType.GetOrZero(),
	}
	return a.sync.GetHttpClient(integration, currentToken, func(integration core.Integration, t *oauth2.Token) error {
		// If token was refreshed we persist the new token info
		_, err = a.db.Auths.Upsert(&models.AuthSetter{
			Provider:     omit.From(string(integration)),
			AccessToken:  omit.From(t.AccessToken),
			RefreshToken: omitnull.From(t.RefreshToken),
			TokenType:    omitnull.From(t.Type()),
			Expiration:   omitnull.From(t.Expiry),
		})
		return err
	})
}
