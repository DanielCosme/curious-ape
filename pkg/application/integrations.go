package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"danicos.dev/daniel/curious-ape/database/gen/models"
	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/integrations/fitbit"
	"danicos.dev/daniel/curious-ape/pkg/integrations/google"
	"danicos.dev/daniel/curious-ape/pkg/oak"
	"danicos.dev/daniel/curious-ape/pkg/persistence"
	"github.com/aarondl/opt/omit"
	"golang.org/x/oauth2"
)

func (a *App) IntegrationsGetList() ([]core.IntegrationInfo, error) {
	var res []core.IntegrationInfo
	for _, integration := range a.sync.IntegrationsList() {
		res = append(res, core.IntegrationInfo{
			Name:   core.ToUpperFist(string(integration)),
			Status: core.IntegrationStatusUnkown,
		})
	}
	return res, nil
}

func (a *App) IntegrationGet(ctx context.Context, provider core.Integration) (core.IntegrationInfo, error) {
	logger := oak.FromContext(ctx).Layer("app")
	defer logger.PopLayer()

	var res core.IntegrationInfo
	var info []string
	var authURL string
	today := core.NewDateToday()
	status := core.IntegrationStatusNotImplemented

	switch provider {
	case core.IntegrationHevy:
		count, err := a.sync.Hevy.Workouts.Count()
		if err != nil {
			status = core.IntegrationStatusDisconnected
			info = append(info, err.Error())
		} else {
			status = core.IntegrationStatusConnected
			info = append(info, fmt.Sprintf("Number of workouts: %d", count))
		}
		res = core.IntegrationInfo{
			Name:   "Hevy",
			Info:   info,
			Status: status,
		}
	case core.IntegrationGoogle:
		_, err := a.fitnessLogsFromGoogle(today)
		if err != nil {
			authURL = a.sync.GenerateOauth2URI(provider)
			if authURL != "" {
				status = core.IntegrationStatusDisconnected
			}
			info = append(info, err.Error())
		} else {
			status = core.IntegrationStatusConnected
		}
		res = core.IntegrationInfo{
			Name:    "Google",
			Info:    info,
			AuthURL: authURL,
			Status:  status,
		}
	case core.IntegrationFitbit:
		sls, err := a.sleepLogsGetFromFitbit(today)
		if err != nil {
			authURL = a.sync.GenerateOauth2URI(provider)
			if authURL != "" {
				status = core.IntegrationStatusDisconnected
			}
			info = append(info, err.Error())
		} else {
			status = core.IntegrationStatusConnected
			if len(sls) > 0 {
				info = append(info, fmt.Sprintf("Total time asleep last night: %s", sls[0].TimeAsleep.String()))
			}
		}
		res = core.IntegrationInfo{
			Name:    "Fitbit",
			Info:    info,
			AuthURL: authURL,
			Status:  status,
		}
	case core.IntegrationToggl:
		profile, err := a.sync.TogglAPI.Me.GetProfile()
		if err != nil {
			info = append(info, err.Error())
		} else if profile != nil {
			status = core.IntegrationStatusConnected
			name := fmt.Sprintf("Profile name: %s", profile.FullName)
			timeZone := fmt.Sprintf("Timezone: %s", profile.Timezone)
			info = append(info, name, timeZone)

			ws, err := a.sync.TogglAPI.Workspace.Get()
			if err == nil {
				for _, w := range ws {
					info = append(info, fmt.Sprintf("Workspace: %s - ID: %d", w.Name, w.ID))
				}
			} else {
				logger.Error(err.Error())
			}

			summary, err := a.sync.TogglAPI.Reports.GetDaySummary(time.Now())
			if err != nil {
				logger.Error(err.Error())
			} else {
				info = append(info, fmt.Sprintf("Total time worked so far: %s", summary.TotalDuration))
			}
		}
		res = core.IntegrationInfo{
			Name:   "Toggl",
			Info:   info,
			Status: status,
		}
	default:
	}
	logger.Debug("integration: "+res.Name, "status", res.Status)

	return res, nil
}

func (a *App) Oauth2Success(provider, code string) error {
	token, err := a.sync.ExchangeToken(core.Integration(provider), code)
	if err != nil {
		return err
	}
	_, err = a.db.Auths.Upsert(&models.OauthTokenSetter{
		Provider:     omit.From(provider),
		AccessToken:  omit.From(token.AccessToken),
		RefreshToken: omit.From(token.RefreshToken),
		TokenType:    omit.From(token.Type()),
		Expiration:   omit.From(token.Expiry),
	})
	a.Log.Info("Authentication successful", "provider", provider, "code", code)
	return err
}

func (a *App) fitbitClient() (res fitbit.API, err error) {
	client, err := a.integrationsGetHttpClient(core.IntegrationFitbit)
	res = fitbit.NewAPI(client)
	return
}

func (a *App) googleClient() (res google.API, err error) {
	client, err := a.integrationsGetHttpClient(core.IntegrationGoogle)
	res = google.NewAPI(client)
	return
}

func (a *App) integrationsGetHttpClient(integration core.Integration) (*http.Client, error) {
	o, err := a.db.Auths.Get(persistence.AuthParams{Provider: integration})
	if err != nil {
		return nil, err
	}
	currentToken := &oauth2.Token{
		AccessToken:  o.AccessToken,
		RefreshToken: o.RefreshToken,
		Expiry:       o.Expiration,
		TokenType:    o.TokenType,
	}
	return a.sync.GetHttpClient(integration, currentToken, func(integration core.Integration, t *oauth2.Token) error {
		// If token was refreshed we persist the new token info
		_, err = a.db.Auths.Upsert(&models.OauthTokenSetter{
			Provider:     omit.From(string(integration)),
			AccessToken:  omit.From(t.AccessToken),
			RefreshToken: omit.From(t.RefreshToken),
			TokenType:    omit.From(t.Type()),
			Expiration:   omit.From(t.Expiry),
		})
		return err
	})
}
