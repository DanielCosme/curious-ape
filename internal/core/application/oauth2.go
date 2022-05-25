package application

import (
	"context"
	"errors"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/core/repository"
	"golang.org/x/oauth2"
	"net/http"
)

func (a *App) Oauth2ConnectProvider(provider string) (string, error) {
	p := entity.IntegrationProvider(provider)
	_, err := a.db.Oauths.Get(entity.Oauth2Filter{
		Providers: []entity.IntegrationProvider{p},
	})
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			if err = a.db.Oauths.Create(&entity.Oauth2{Provider: p}); err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	config := a.oauth2GetConfigurationForProvider(p)
	uri := config.AuthCodeURL("")
	return uri, nil
}

func (a *App) Oauth2Success(provider, code string) error {
	p := entity.IntegrationProvider(provider)
	config := a.oauth2GetConfigurationForProvider(p)

	t, err := config.Exchange(context.Background(), code)
	if err != nil {
		return err
	}

	o, err := a.db.Oauths.Get(entity.Oauth2Filter{Providers: []entity.IntegrationProvider{p}})
	if err != nil {
		return err
	}

	o.AccessToken = t.AccessToken
	o.RefreshToken = t.RefreshToken
	o.Expiration = t.Expiry
	o.Type = t.Type()

	_, err = a.db.Oauths.Update(o)
	return err
}

func (a *App) Oauth2GetClient(provider entity.IntegrationProvider) (*http.Client, error) {
	config := a.oauth2GetConfigurationForProvider(provider)

	o, err := a.db.Oauths.Get(entity.Oauth2Filter{Provider: provider})
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{
		AccessToken:  o.AccessToken,
		RefreshToken: o.RefreshToken,
		Expiry:       o.Expiration,
		TokenType:    o.Type,
	}

	// Check if token is still valid, if not refresh it
	newToken, err := config.TokenSource(context.Background(), t).Token()
	if newToken.AccessToken != t.AccessToken {
		// If token was refreshed we persist the new token info
		_, err = a.db.Oauths.Update(&entity.Oauth2{
			ID:           o.ID,
			AccessToken:  newToken.AccessToken,
			RefreshToken: newToken.RefreshToken,
			Type:         newToken.TokenType,
			Expiration:   newToken.Expiry,
		})
	}

	return config.Client(context.Background(), newToken), err
}

func (a *App) oauth2GetConfigurationForProvider(provider entity.IntegrationProvider) *oauth2.Config {
	var config *entity.Oauth2Config
	switch provider {
	case entity.ProviderFitbit:
		config = a.env.Fitbit
	default:
		return nil
	}

	return &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.AuthURL,
			TokenURL: config.TokenURL,
		},
		RedirectURL: config.RedirectURL,
		Scopes:      config.Scopes,
	}
}
