package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database"
	"github.com/danielcosme/curious-ape/internal/database/gen/models"
	"golang.org/x/oauth2"
	"time"
)

type IntegrationInfo struct {
	Name        string
	State       IntegrationState
	ProfileInfo string
	AuthURL     string
	Problem     string
}

func (a *App) IntegrationsGet() ([]IntegrationInfo, error) {
	// TODO Improve this.
	var integrations []IntegrationInfo
	// I have to determine the status beforehand.

	state := IntegrationDisconnected
	var profileInfo, authURL, problem string
	o, err := a.db.Auths.Get(database.AuthParams{Provider: core.IntegrationFitbit})
	if database.IfNotFoundErr(err) {
		return nil, err
	}
	if o == nil {
		authURL = a.generateOauth2URI(core.IntegrationFitbit)
	} else {
		currentToken := &oauth2.Token{
			AccessToken:  o.AccessToken,
			RefreshToken: o.RefreshToken.GetOrZero(),
			Expiry:       o.Expiration.GetOrZero(),
			TokenType:    o.TokenType.GetOrZero(),
		}
		// Refresh token.
		newToken, err := a.cfg.Fitbit.TokenSource(context.Background(), currentToken).Token()
		if err != nil {
			return nil, err
		}
		// Check if token is still valid.
		if newToken.AccessToken != currentToken.AccessToken {
			// If token was refreshed we persist the new token info
			_, err = a.db.Auths.Upsert(&models.AuthSetter{
				Provider:     omit.From(core.IntegrationFitbit),
				AccessToken:  omit.From(newToken.AccessToken),
				RefreshToken: omitnull.From(newToken.RefreshToken),
				TokenType:    omitnull.From(newToken.Type()),
				Expiration:   omitnull.From(newToken.Expiry),
			})
			if err != nil {
				return nil, err
			}
			currentToken = newToken
		}
		fitbitHttp := a.cfg.Fitbit.Client(context.Background(), currentToken)
		fitbitAPI := a.sync.FitbitClient(fitbitHttp)
		sleepRecord, err := fitbitAPI.Sleep.GetByDate(time.Now())
		if err != nil {
			problem = err.Error()
			authURL = a.generateOauth2URI(core.IntegrationFitbit)
		} else {
			state = IntegrationConnected
			profileInfo = fmt.Sprintf(
				"Today's total time in bed is: %s",
				ToDuration(sleepRecord.Summary.TotalTimeInBed).String())
		}
	}

	fitbitStatus := IntegrationInfo{
		Name:        "Fitbit",
		State:       state,
		ProfileInfo: profileInfo,
		AuthURL:     authURL,
		Problem:     problem,
	}
	integrations = append(integrations, fitbitStatus)

	return integrations, nil
}

func ToDuration(i int) time.Duration {
	return time.Duration(i) * time.Minute
}

func (a *App) generateOauth2URI(provider core.Integration) string {
	var opts []oauth2.AuthCodeOption
	var config *oauth2.Config

	switch provider {
	case core.IntegrationFitbit:
		config = a.cfg.Fitbit
	case core.IntegrationGoogle:
		config = a.cfg.Google
		opts = append(opts,
			oauth2.SetAuthURLParam("access_type", "offline"),
			oauth2.SetAuthURLParam("approval_prompt", "force"),
		)
	default:
		return ""
	}

	return config.AuthCodeURL("", opts...)
}

func (a *App) Oauth2Success(provider, code string) error {
	var config *oauth2.Config

	switch provider {
	case core.IntegrationFitbit:
		config = a.cfg.Fitbit
	case core.IntegrationGoogle:
		config = a.cfg.Google
	default:
		return errors.New("non-implemented provider: " + string(provider))
	}

	token, err := config.Exchange(context.Background(), code)
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
	a.Log.Info("Exchange successful", "token_type", token.TokenType)
	return err
}
