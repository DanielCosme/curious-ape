package application

import (
	"fmt"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/database"
	"github.com/danielcosme/curious-ape/pkg/database/gen/models"
	"github.com/danielcosme/curious-ape/pkg/integrations/fitbit"
	"github.com/danielcosme/curious-ape/pkg/integrations/google"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

type IntegrationInfo struct {
	Name        string `json:"name"`
	IsConnected bool   `json:"isConnected"`
	Info        []string
	AuthURL     string
	Problem     string
}

func (a *App) IntegrationsGetAll() ([]IntegrationInfo, error) {
	var res []IntegrationInfo
	for _, integration := range a.sync.IntegrationsList() {
		res = append(res, IntegrationInfo{
			Name:        core.ToUpperFist(string(integration)),
			IsConnected: false,
		})
	}
	return res, nil
}

func (a *App) IntegrationsGet(provider core.Integration) (IntegrationInfo, error) {
	var res IntegrationInfo
	var info []string
	var authURL, problem string
	isConnected := false
	today := core.NewDateToday()

	switch provider {
	case core.IntegrationGoogle:
		_, err := a.fitnessLogsFromGoogle(today)
		if err != nil {
			authURL = a.sync.GenerateOauth2URI(provider)
			problem = err.Error()
		} else {
			isConnected = true
		}
		res = IntegrationInfo{
			Name:        "Google",
			IsConnected: isConnected,
			Info:        info,
			AuthURL:     authURL,
			Problem:     problem,
		}
	case core.IntegrationFitbit:
		sls, err := a.sleepLogsGetFromFitbit(today)
		if err != nil {
			authURL = a.sync.GenerateOauth2URI(provider)
			problem = err.Error()
		} else {
			isConnected = true
			if len(sls) > 0 {
				minutes := sls[0].MinutesAsleep.GetOrZero()
				dur := time.Duration(minutes) * time.Minute
				info = append(info, fmt.Sprintf("Total time asleep last night: %s", dur.String()))
			}
		}
		res = IntegrationInfo{
			Name:        "Fitbit",
			IsConnected: isConnected,
			Info:        info,
			AuthURL:     authURL,
			Problem:     problem,
		}
	case core.IntegrationToggl:
		profile, err := a.sync.TogglAPI.Me.GetProfile()
		if err != nil {
			problem = err.Error()
		} else if profile != nil {
			isConnected = true
			name := fmt.Sprintf("Profile name: %s", profile.FullName)
			timeZone := fmt.Sprintf("Timezone: %s", profile.Timezone)
			info = append(info, name, timeZone)

			ws, err := a.sync.TogglAPI.Workspace.Get()
			if err == nil {
				for _, w := range ws {
					info = append(info, fmt.Sprintf("Workspace: %s - ID: %d", w.Name, w.ID))
				}
			} else {
				a.Log.Error(err.Error())
			}

			summary, err := a.sync.TogglAPI.Reports.GetDaySummary(time.Now())
			if err != nil {
				a.Log.Error(err.Error())
			} else {
				info = append(info, fmt.Sprintf("Total time so far: %s", summary.TotalDuration))
			}
		}
		res = IntegrationInfo{
			Name:        "Toggl",
			IsConnected: isConnected,
			Info:        info,
			AuthURL:     authURL,
			Problem:     problem,
		}
	default:
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
