package application

import (
	"net/http"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/gen/bob/models"
	"danicos.dev/daniel/curious-ape/pkg/integrations/fitbit"
	"danicos.dev/daniel/curious-ape/pkg/integrations/google"
	"danicos.dev/daniel/curious-ape/pkg/persistence"
	"github.com/aarondl/opt/omit"
	"golang.org/x/oauth2"
)

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
