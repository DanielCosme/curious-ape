package application

import (
	"context"
	"errors"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/core/repository"
	"golang.org/x/oauth2"
)

func (a *App) Oauth2ConnectProvider(provider string) (string, error) {
	p := entity.IntegrationProviders(provider)
	_, err := a.db.Oauths.Get(entity.Oauth2Filter{
		Providers: []entity.IntegrationProviders{p},
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

	config := a.Oauth2GetConfigurationForProvider(p)
	uri := config.AuthCodeURL("")
	return uri, nil
}

func (a *App) Oauth2Success(provider, code string) error {
	p := entity.IntegrationProviders(provider)
	config := a.Oauth2GetConfigurationForProvider(p)

	t, err := config.Exchange(context.Background(), code)
	if err != nil {
		return err
	}

	o, err := a.db.Oauths.Get(entity.Oauth2Filter{Providers: []entity.IntegrationProviders{p}})
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

func (a *App) Oauth2GetConfigurationForProvider(provider entity.IntegrationProviders) *oauth2.Config {
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
