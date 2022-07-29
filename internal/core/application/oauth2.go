package application

import (
	"context"
	"github.com/danielcosme/curious-ape/internal/core/database"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/sdk/errors"
	"golang.org/x/oauth2"
	"net/http"
	"strconv"
)

func (a *App) Oauth2ConnectProvider(provider string) (string, error) {
	p := entity.IntegrationProvider(provider)
	o, err := a.db.Oauths.Get(entity.Oauth2Filter{Provider: []entity.IntegrationProvider{p}})
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return "", err
	}
	if o == nil {
		if err = a.db.Oauths.Create(&entity.Oauth2{Provider: p}); err != nil {
			return "", err
		}
	}

	config := a.oauth2GetConfigurationForProvider(p)

	opts := []oauth2.AuthCodeOption{}
	switch p {
	case entity.ProviderGoogle:
		opts = append(opts,
			oauth2.SetAuthURLParam("access_type", "offline"),
			oauth2.SetAuthURLParam("approval_prompt", "force"),
		)
	}

	uri := config.AuthCodeURL("", opts...)
	return uri, nil
}

func (a *App) Oauth2Success(provider, code string) error {
	p := entity.IntegrationProvider(provider)
	config := a.oauth2GetConfigurationForProvider(p)

	t, err := config.Exchange(context.Background(), code)
	if err != nil {
		return err
	}

	o, err := a.db.Oauths.Get(entity.Oauth2Filter{Provider: []entity.IntegrationProvider{p}})
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

	o, err := a.db.Oauths.Get(entity.Oauth2Filter{Provider: []entity.IntegrationProvider{provider}})
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{
		AccessToken:  o.AccessToken,
		RefreshToken: o.RefreshToken,
		Expiry:       o.Expiration,
		TokenType:    o.Type,
	}

	switch provider {
	case entity.ProviderFitbit:
		newToken, err := config.TokenSource(context.Background(), token).Token()
		// Check if token is still valid, if not refresh it
		if newToken.AccessToken != token.AccessToken {
			// If token was refreshed we persist the new token info
			_, err = a.db.Oauths.Update(&entity.Oauth2{
				ID:           o.ID,
				AccessToken:  newToken.AccessToken,
				RefreshToken: newToken.RefreshToken,
				Type:         newToken.TokenType,
				Expiration:   newToken.Expiry,
			})
			token = newToken
		}
		if err != nil {
			return nil, err
		}
	default:
		// by default, we assume these tokens do not expire
	}

	return config.Client(context.Background(), token), err
}

func (a *App) oauth2GetConfigurationForProvider(provider entity.IntegrationProvider) *oauth2.Config {
	var config *entity.Oauth2Config
	switch provider {
	case entity.ProviderFitbit:
		config = a.cfg.Fitbit
	case entity.ProviderGoogle:
		config = a.cfg.Google
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

func (a *App) Oauth2AddToken(token, provider string) (string, error) {
	if token == "" {
		return "", errors.New("token is empty")
	}

	switch entity.IntegrationProvider(provider) {
	case entity.ProviderToggl:
		o, err := a.db.Oauths.Get(entity.Oauth2Filter{Provider: []entity.IntegrationProvider{entity.ProviderToggl}})
		if err != nil && !errors.Is(err, database.ErrNotFound) {
			return "", err
		}

		if o == nil {
			o = &entity.Oauth2{
				Provider:    entity.ProviderToggl,
				AccessToken: token,
				Type:        "Bearer",
			}
			if err := a.db.Oauths.Create(o); err != nil {
				return "", err
			}
		} else {
			o.AccessToken = token
			if _, err := a.db.Oauths.Update(o); err != nil {
				return "", err
			}
		}
		api := a.Sync.TogglClient(o.AccessToken)
		me, err := api.Me.GetProfile()
		if err != nil {
			return "", err
		}

		ws, err := api.Workspace.Get()
		if err != nil {
			a.Log.Error(err)
		}
		if len(ws) != 1 {
			return "", errors.New("only one workspace is supported for toggl")
		}
		w := ws[0]
		o.ToogglOrganizationID = strconv.Itoa(w.OrganizationID)
		o.ToogglWorkSpaceID = strconv.Itoa(w.ID)
		if _, err := a.db.Oauths.Update(o); err != nil {
			return "", err
		}

		return me.Fullname, nil
	default:
		return "", errors.New("invalid provider")
	}

	return "", nil
}
