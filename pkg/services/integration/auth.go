package integration

import (
	"log/slog"
	"net/http"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/gen/bob/models"
	"danicos.dev/daniel/curious-ape/pkg/integrations/fitbit"
	"danicos.dev/daniel/curious-ape/pkg/integrations/google"
	"github.com/aarondl/opt/omit"
	"golang.org/x/oauth2"
)

func (svc *Service) Oauth2Success(provider, code string) error {
	token, err := svc.sync.ExchangeToken(core.Integration(provider), code)
	if err != nil {
		return err
	}
	_, err = upsert(svc.db, &models.OauthTokenSetter{
		Provider:     omit.From(provider),
		AccessToken:  omit.From(token.AccessToken),
		RefreshToken: omit.From(token.RefreshToken),
		TokenType:    omit.From(token.Type()),
		Expiration:   omit.From(token.Expiry),
	})
	slog.Info("Authentication successful", "provider", provider, "code", code)
	return err
}

func (svc *Service) fitbitClient() (res fitbit.API, err error) {
	client, err := svc.integrationsGetHttpClient(core.IntegrationFitbit)
	res = fitbit.NewAPI(client)
	return
}

func (svc *Service) googleClient() (res google.API, err error) {
	client, err := svc.integrationsGetHttpClient(core.IntegrationGoogle)
	res = google.NewAPI(client)
	return
}

func (svc *Service) integrationsGetHttpClient(integration core.Integration) (*http.Client, error) {
	o, err := get(svc.db, AuthParams{Provider: integration})
	if err != nil {
		return nil, err
	}
	currentToken := &oauth2.Token{
		AccessToken:  o.AccessToken,
		RefreshToken: o.RefreshToken,
		Expiry:       o.Expiration,
		TokenType:    o.TokenType,
	}
	return svc.sync.GetHttpClient(integration, currentToken, func(integration core.Integration, t *oauth2.Token) error {
		// If token was refreshed we persist the new token info
		_, err = upsert(svc.db, &models.OauthTokenSetter{
			Provider:     omit.From(string(integration)),
			AccessToken:  omit.From(t.AccessToken),
			RefreshToken: omit.From(t.RefreshToken),
			TokenType:    omit.From(t.Type()),
			Expiration:   omit.From(t.Expiry),
		})
		return err
	})
}
